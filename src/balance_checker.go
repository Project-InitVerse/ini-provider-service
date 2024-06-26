package provider

import (
	"context"
	"math/rand"
	"time"

	"github.com/boz/go-lifecycle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	btypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/libs/log"
	tmrpc "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ovrclk/akash/pubsub"
	"github.com/ovrclk/akash/util/runner"
	dtypes "github.com/ovrclk/akash/x/deployment/types/v1beta2"
	"github.com/ovrclk/akash/x/escrow/client/util"
	mtypes "github.com/ovrclk/akash/x/market/types/v1beta2"

	"providerService/src/event"
	"providerService/src/session"

	aclient "github.com/ovrclk/akash/client"
	netutil "github.com/ovrclk/akash/util/network"
)

type respState int

const (
	respStateNextCheck = iota
	respStateOutOfFunds
	respStateScheduledWithdraw
)

// BalanceCheckerConfig is struct
type BalanceCheckerConfig struct {
	WithdrawalPeriod        time.Duration
	LeaseFundsCheckInterval time.Duration
}

type leaseState struct {
	tm                  *time.Timer
	scheduledWithdrawAt time.Time
}

type balanceChecker struct {
	ctx     context.Context
	session session.Session
	log     log.Logger
	lc      lifecycle.Lifecycle
	bus     pubsub.Bus
	ownAddr sdk.AccAddress
	bqc     btypes.QueryClient
	aqc     aclient.QueryClient
	leases  map[mtypes.LeaseID]*leaseState
	cfg     BalanceCheckerConfig
}

type leaseCheckResponse struct {
	lid        mtypes.LeaseID
	checkAfter time.Duration
	state      respState
	err        error
}

func (bc *balanceChecker) runEscrowCheck(ctx context.Context, lid mtypes.LeaseID, scheduledWithdraw bool, res chan<- leaseCheckResponse) {
	go func() {
		select {
		case <-bc.lc.Done():
		case res <- bc.doEscrowCheck(ctx, lid, scheduledWithdraw):
		}
	}()
}

func (bc *balanceChecker) doEscrowCheck(ctx context.Context, lid mtypes.LeaseID, scheduledWithdraw bool) leaseCheckResponse {
	resp := leaseCheckResponse{
		lid:   lid,
		state: respStateNextCheck,
	}

	if scheduledWithdraw {
		resp.state = respStateScheduledWithdraw
	}

	var syncInfo *tmrpc.SyncInfo
	syncInfo, resp.err = bc.session.Client().NodeSyncInfo(ctx)
	if resp.err != nil {
		return resp
	}

	if syncInfo.CatchingUp {
		resp.err = aclient.ErrNodeNotSynced
		return resp
	}

	var dResp *dtypes.QueryDeploymentResponse
	var lResp *mtypes.QueryLeasesResponse

	// Fetch the balance of the escrow account
	dResp, resp.err = bc.aqc.Deployment(ctx, &dtypes.QueryDeploymentRequest{
		ID: lid.DeploymentID(),
	})

	if resp.err != nil {
		return resp
	}

	lResp, resp.err = bc.aqc.Leases(ctx, &mtypes.QueryLeasesRequest{
		Filters: mtypes.LeaseFilters{
			Owner: lid.Owner,
			DSeq:  lid.DSeq,
			State: "active",
		},
	})

	if resp.err != nil {
		return resp
	}

	totalLeaseAmount := sdk.NewDec(0)
	for _, lease := range lResp.Leases {
		totalLeaseAmount = totalLeaseAmount.Add(lease.Lease.Price.Amount)
	}

	balanceRemain := util.LeaseCalcBalanceRemain(dResp.EscrowAccount.TotalBalance().Amount,
		syncInfo.LatestBlockHeight,
		dResp.EscrowAccount.SettledAt,
		totalLeaseAmount)

	blocksRemain := util.LeaseCalcBlocksRemain(balanceRemain, totalLeaseAmount)

	// lease is out of funds
	if blocksRemain <= 0 {
		resp.state = respStateOutOfFunds
		resp.checkAfter = time.Minute * 10
	} else {
		blocksPerCheckInterval := int64(bc.cfg.LeaseFundsCheckInterval / netutil.AverageBlockTime)
		if blocksRemain > blocksPerCheckInterval {
			blocksRemain = blocksPerCheckInterval
		}

		resp.checkAfter = time.Duration(blocksRemain) * netutil.AverageBlockTime
	}

	return resp
}

func (bc *balanceChecker) startWithdraw(lid mtypes.LeaseID) error {
	msg := &mtypes.MsgWithdrawLease{
		LeaseID: lid,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errch := make(chan error, 1)

	go func(ch chan<- error) {
		ch <- bc.session.Client().Tx().Broadcast(ctx, msg)
	}(errch)

	select {
	case <-bc.lc.Done():
		// give request extra 30s to finish before force canceling
		select {
		case <-time.After(30 * time.Second):
			cancel()
			return bc.lc.Error()
		case err := <-errch:
			return err
		}
	case err := <-errch:
		return err
	}
}

func (bc *balanceChecker) run(startCh chan<- error) {
	ctx, cancel := context.WithCancel(bc.ctx)

	defer func() {
		cancel()
		bc.lc.ShutdownCompleted()

		for _, lState := range bc.leases {
			if lState.tm != nil && !lState.tm.Stop() {
				<-lState.tm.C
			}
		}

		bc.log.Debug("shutdown complete")
	}()

	leaseCheckCh := make(chan leaseCheckResponse, 1)
	var resultch chan runner.Result

	subscriber, err := bc.bus.Subscribe()
	startCh <- err
	if err != nil {
		return
	}

	resultch = make(chan runner.Result, 1)

loop:
	for {
		select {
		case <-bc.lc.ShutdownRequest():
			bc.log.Debug("shutting down")
			bc.lc.ShutdownInitiated(nil)
			break loop
		case evt := <-subscriber.Events():
			switch ev := evt.(type) {
			case event.LeaseAddFundsMonitor:
				var scheduledWithdraw time.Time
				// if provider configured with periodic force withdrawal
				// set next time at which withdraw will happen
				if bc.cfg.WithdrawalPeriod > 0 {
					scheduledWithdraw = time.Now().Add(bc.cfg.WithdrawalPeriod)
				}

				lState := &leaseState{
					scheduledWithdrawAt: scheduledWithdraw,
				}

				bc.leases[ev.LeaseID] = lState

				// if there was provider restart with a bunch of active leases
				// spread their requests across 1min interval
				// to reduce pressure on the RPC
				if !ev.IsNewLease {
					checkIn := time.Duration(rand.Int63n(int64(time.Minute))) // nolint: gosec
					lState.tm = bc.timerFunc(ctx, checkIn, ev.LeaseID, false, leaseCheckCh)
				} else {
					bc.runEscrowCheck(ctx, ev.LeaseID, false, leaseCheckCh)
				}
			case event.LeaseRemoveFundsMonitor:
				lsState, exists := bc.leases[ev.LeaseID]
				if !exists {
					break
				}

				if lsState.tm != nil && !lsState.tm.Stop() {
					<-lsState.tm.C
				}

				delete(bc.leases, ev.LeaseID)
			}
		case res := <-leaseCheckCh:
			// we may have timer fired just a heart beat ahead of lease remove event.
			lState, exists := bc.leases[res.lid]
			if !exists {
				continue loop
			}

			withdraw := false

			switch res.state {
			case respStateOutOfFunds:
				bc.log.Debug("lease is out of funds", "lease", res.lid)
				// reschedule funds check. if lease not being topped up then network will close it
				fallthrough
			case respStateScheduledWithdraw:
				withdraw = true
				bc.log.Debug("sending withdraw", "lease", res.lid)
				// reschedule periodic withdraw if configured
				if bc.cfg.WithdrawalPeriod > 0 {
					lState.scheduledWithdrawAt = time.Now().Add(bc.cfg.WithdrawalPeriod)
				}

				fallthrough
			case respStateNextCheck:
				timerPeriod := res.checkAfter
				scheduledWithdraw := false

				if res.err != nil {
					bc.log.Info("couldn't check lease balance. retrying in 1m", "leaseId", res.lid, "error", res.err.Error())
					timerPeriod = time.Minute
				} else if !withdraw && !lState.scheduledWithdrawAt.IsZero() {
					withdrawIn := time.Until(lState.scheduledWithdrawAt)
					if timerPeriod >= withdrawIn {
						timerPeriod = withdrawIn
						scheduledWithdraw = true
					}
				}

				lState.tm = bc.timerFunc(ctx, timerPeriod, res.lid, scheduledWithdraw, leaseCheckCh)
			}

			if withdraw {
				go func() {
					select {
					case <-bc.ctx.Done():
					case resultch <- runner.NewResult(res.lid, bc.startWithdraw(res.lid)):
					}
				}()
			}
		case res := <-resultch:
			if err := res.Error(); err != nil {
				bc.log.Error("failed to do lease withdrawal", "err", err, "LeaseID", res.Value().(mtypes.LeaseID))
			}
		}
	}
}

func (bc *balanceChecker) timerFunc(ctx context.Context, d time.Duration, lid mtypes.LeaseID, scheduledWithdraw bool, ch chan<- leaseCheckResponse) *time.Timer {
	return time.AfterFunc(d, func() {
		select {
		case <-bc.ctx.Done():
		case ch <- bc.doEscrowCheck(ctx, lid, scheduledWithdraw):
		}
	})
}
