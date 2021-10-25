package debug

import (
	"time"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
)

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested in
// depth-first post-order.
//
// args contains the arguments for GenerateRules. This is passed as a struct to
// avoid breaking implementations in the future when new fields are added.
//
// A GenerateResult struct is returned. Optional fields may be added to this
// type in the future.
//
// Any non-fatal errors this function encounters should be logged using
// log.Print.
func (dl *debugLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	dc := getDebugConfig(args.Config)

	for _, f := range args.RegularFiles {
		dc.Debug().Str("file", f).Msg("read dir")
	}

	for _, r := range args.OtherGen {
		dc.Debug().
			Str("name", r.Name()).
			Str("kind", r.Kind()).
			Stringer("label", label.New("", args.Rel, r.Name())).
			Msg("generated rule")
	}

	dc.increment()
	// dc.done(args.Rel)

	current := time.Now()
	diff := current.Sub(dl.prev)
	elapsed := current.Sub(dl.start)
	dl.prev = current

	if diff.Milliseconds() > 1000 {
		dc.Warn().
			Str("pkg", label.New("", args.Rel, "all").String()).
			Stringer("t", diff).
			Int("rule-count", len(args.OtherGen)).
			Int("file-count", len(args.RegularFiles)).
			Msgf("%s: slow (%s)", elapsed.Round(time.Second), diff.Round(time.Second))
	}

	if args.Rel == "" {
		go func() {
			time.Sleep(0)
			dc.Error().Msg("")
			dc.Error().Msg(`
This is taking longer than expected...

A common cause is gazelle failing to locally resolve a golang package.
For example, an importpath like 'github.com/robinhoodmarkets/rh/i/dont/exist will 
trigger gazelle to clone the monorepo again, only to discover 
`)
		}()
	}

	return language.GenerateResult{}
}
