package debug

import (
	"flag"
	"fmt"
	"path"
	"time"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/rs/zerolog"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

// The following methods are implemented to satisfy the
// https://pkg.go.dev/github.com/bazelbuild/bazel-gazelle/resolve?tab=doc#Resolver
// interface, but are otherwise unused.
func (dl *debugLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	dc := newDebugConfig()
	c.Exts[DebugLangName] = dc
	dl.start = time.Now()
	dl.prev = time.Now()
	dl.progress = mpb.New(mpb.WithWidth(64))
}

func (*debugLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	dc := getDebugConfig(c)

	fs.VisitAll(func(f *flag.Flag) {
		dc.Info().
			Str("name", f.Name).
			Str("value", f.Value.String()).
			Str("default", f.DefValue).
			Msg("checking flag")
	})

	return nil
}

func (*debugLang) KnownDirectives() []string { return []string{"log_level", "progress"} }

// Configure implements config.Configurer
func (dl *debugLang) Configure(c *config.Config, rel string, f *rule.File) {
	var dc *debugConfig
	if raw, ok := c.Exts[DebugLangName]; !ok {
		dc = newDebugConfig()
	} else {
		dc = raw.(*debugConfig).clone()
	}
	c.Exts[DebugLangName] = dc

	dc.Debug().Str("rel", rel).Msg("visiting")

	if f != nil {
		for _, d := range f.Directives {
			dc.Debug().
				Str("key", d.Key).
				Str("value", d.Value).
				Str("rel", rel).
				Msg("configuring directive")
			switch d.Key {
			case "log_level":
				level, err := zerolog.ParseLevel(d.Value)
				if err != nil {
					fmt.Printf("warning: bad log_level: %v", err)
				} else {
					dc.Logger = dc.Logger.Level(level)
				}
			case "progress":
				dirCount := dc.countPackages(path.Join(c.WorkDir, rel))
				dc.Warn().Int64("count", dirCount).Str("work_dir", c.WorkDir).Msg("total dirs")

				dc.wantProgress = true
				dc.bar = dl.progress.Add(dirCount,
					// progress bar filler with customized style
					mpb.NewBarFiller(mpb.BarStyle().Lbound("╢").Filler("▌").Tip("▌").Padding("░").Rbound("╟")),
					mpb.PrependDecorators(
						decor.Name(path.Base(rel)),
						decor.Percentage(decor.WCSyncSpace),
					),
					mpb.AppendDecorators(
						decor.OnComplete(
							decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done",
						),
					),
					mpb.BarRemoveOnComplete(),
				)

			case "exclude":
				if err := checkPathMatchPattern(path.Join(rel, d.Value)); err != nil {
					dc.Warn().Str("rel", rel).Str("pattern", d.Value).Err(err).Msg("exclusion pattern not valid")
					continue
				}
				dc.excludes = append(dc.excludes, path.Join(rel, d.Value))
			}
		}
	}

}
