package util

import (
	"os"

	"github.com/go-kit/kit/log/term" // nolint: staticcheck
	"github.com/tendermint/tendermint/libs/log"
)

// OpenLogger is function get logger
func OpenLogger() log.Logger {
	// logger with no color output - current debug colors are invisible for me.
	return log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), func(_ ...interface{}) term.FgBgColor {
		return term.FgBgColor{}
	})
}
