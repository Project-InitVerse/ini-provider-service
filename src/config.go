package provider

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/ovrclk/akash/types/v1beta2"
	"github.com/ovrclk/akash/validation/constants"
	mtypes "github.com/ovrclk/akash/x/market/types/v1beta2"

	"github.com/ovrclk/provider-services/bidengine"
)

// Config is struct
type Config struct {
	ClusterWaitReadyDuration        time.Duration
	ClusterPublicHostname           string
	ClusterExternalPortQuantity     uint
	InventoryResourcePollPeriod     time.Duration
	InventoryResourceDebugFrequency uint
	BidPricingStrategy              bidengine.BidPricingStrategy
	BidDeposit                      sdk.Coin
	CPUCommitLevel                  float64
	MemoryCommitLevel               float64
	StorageCommitLevel              float64
	MaxGroupVolumes                 int
	BlockedHostnames                []string
	BidTimeout                      time.Duration
	ManifestTimeout                 time.Duration
	BalanceCheckerCfg               BalanceCheckerConfig
	Attributes                      types.Attributes
	DeploymentIngressStaticHosts    bool
	DeploymentIngressDomain         string
	ClusterSettings                 map[interface{}]interface{}
	RPCQueryTimeout                 time.Duration
	CachedResultMaxAge              time.Duration
}

// NewDefaultConfig create default config
func NewDefaultConfig() Config {
	return Config{
		ClusterWaitReadyDuration: time.Second * 10,
		BidDeposit:               mtypes.DefaultBidMinDeposit,
		BalanceCheckerCfg: BalanceCheckerConfig{
			LeaseFundsCheckInterval: 1 * time.Minute,
			WithdrawalPeriod:        24 * time.Hour,
		},
		MaxGroupVolumes: constants.DefaultMaxGroupVolumes,
	}
}
