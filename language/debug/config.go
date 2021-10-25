package debug

import (
	"flag"
	"fmt"
	"time"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/rs/zerolog"
)

// The following methods are implemented to satisfy the
// https://pkg.go.dev/github.com/bazelbuild/bazel-gazelle/resolve?tab=doc#Resolver
// interface, but are otherwise unused.
func (dl *debugLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	dc := newDebugConfig()
	c.Exts[DebugLangName] = dc

	dl.start = time.Now()
	dl.prev = time.Now()
}

func (*debugLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	dc := getDebugConfig(c)

	fs.VisitAll(func(f *flag.Flag) {
		dc.Debug().
			Str("name", f.Name).
			Str("value", f.Value.String()).
			Str("default", f.DefValue).
			Msg("checking flag")
	})

	return nil
}

func (*debugLang) KnownDirectives() []string {
	return []string{"log_level", "generaterules_slow_warn_duration", "show_total_elapsed_time_info_messages"}
}

// Configure implements config.Configurer
func (dl *debugLang) Configure(c *config.Config, rel string, f *rule.File) {
	var dc *debugConfig
	if raw, ok := c.Exts[DebugLangName]; !ok {
		dc = newDebugConfig()
	} else {
		dc = raw.(*debugConfig).clone()
	}
	c.Exts[DebugLangName] = dc

	dc.Debug().Str("dir", rel).Msg("visiting")

	if f != nil {
		for _, d := range f.Directives {
			dc.Debug().
				Str("key", d.Key).
				Str("value", d.Value).
				Str("dir", rel).
				Msg("configuring directive")
			switch d.Key {
			case "log_level":
				level, err := zerolog.ParseLevel(d.Value)
				if err != nil {
					fmt.Printf("warning: bad log_level: %v", err)
				} else {
					dc.Logger = dc.Logger.Level(level)
				}
			case "show_total_elapsed_time_info_messages":
				switch d.Value {
				case "true":
					dc.showTotalElapsedTimeMessages = true
				case "false":
					dc.showTotalElapsedTimeMessages = false
				}
			case "generaterules_slow_warn_duration":
				threshold, err := time.ParseDuration(d.Value)
				if err != nil {
					fmt.Printf("warning: bad generaterules_slow_warn_duration: %v", err)
				} else {
					dc.generaterulesSlowWarnDuration = threshold
				}
			}
		}
	}

}
