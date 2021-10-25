package debug

import (
	"os"
	"path"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bmatcuk/doublestar"
	"github.com/karrick/godirwalk"
	"github.com/rs/zerolog"
	"github.com/vbauerster/mpb/v7"
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
		excludes: []string{"**/.git"},
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
	wantProgress bool
	parent       *debugConfig
	zerolog.Logger
	excludes []string
	bar      *mpb.Bar
}

func (c *debugConfig) clone() *debugConfig {
	return &debugConfig{
		wantProgress: false, // progress bar not inherited
		parent:       c,
		Logger:       c.Logger.With().Logger(),
		excludes:     c.excludes,
	}
}

func (dc *debugConfig) increment() {
	if dc.bar != nil {
		dc.bar.Increment()
	}
	if dc.parent != nil {
		dc.parent.increment()
	}
}

func (dc *debugConfig) done(rel string) {
	if dc.bar != nil {
		dc.bar.Abort(true)
	}
	if dc.parent != nil {
		dc.parent.done(rel)
	}
}

func (dc *debugConfig) isExcluded(osPathname string) bool {
	for _, x := range dc.excludes {
		matched, err := doublestar.Match(x, osPathname)
		if err != nil {
			// doublestar.Match returns only one possible error, and only if the
			// pattern is not valid. During the configuration of the walker (see
			// Configure below), we discard any invalid pattern and thus an error
			// here should not be possible.
			dc.Panic().Msgf("error during doublestar.Match. This should not happen, please file an issue https://github.com/bazelbuild/bazel-gazelle/issues/new: %s", err)
		}
		if matched {
			return true
		}
	}
	return false
}

func (dc *debugConfig) countPackages(rootDir string) int64 {
	var count int64
	godirwalk.Walk(rootDir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsSymlink() {
				return godirwalk.SkipThis
			}
			if path.Base(osPathname) == ".git" {
				return godirwalk.SkipThis
			}
			if dc.isExcluded(osPathname) {
				return godirwalk.SkipThis
			}
			base := path.Base(osPathname)
			if base == "BUILD.bazel" || base == "BUILD" {
				count++
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})
	return count
}

func checkPathMatchPattern(pattern string) error {
	_, err := doublestar.Match(pattern, "x")
	return err
}
