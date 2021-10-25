package debug

import (
	"os"
	"time"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/rs/zerolog"
)

// getDebugConfig returns the debug language configuration. If the debug
// extension was not run, it will return nil.
func getDebugConfig(c *config.Config) *debugConfig {
	dc := c.Exts[DebugLangName]
	if dc == nil {
		return nil
	}
	return dc.(*debugConfig)
}

func newDebugConfig() *debugConfig {
	return &debugConfig{
		Logger: zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			Level(zerolog.WarnLevel).
			With().
			Timestamp().
			Str("lang", DebugLangName).
			Logger(),
	}
}

// debugConfig is the config implementation which embeds a logger
type debugConfig struct {
	zerolog.Logger

	generaterulesSlowWarnDuration time.Duration
	showTotalElapsedTimeMessages  bool
}

func (c *debugConfig) clone() *debugConfig {
	return &debugConfig{
		Logger:                        c.Logger.With().Logger(),
		generaterulesSlowWarnDuration: c.generaterulesSlowWarnDuration,
		showTotalElapsedTimeMessages:  c.showTotalElapsedTimeMessages,
	}
}
