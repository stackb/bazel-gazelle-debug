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
		dc.Debug().
			Str("file", f).
			Str("dir", args.Rel).
			Msg("read dir")
	}

	for _, r := range args.OtherGen {
		dc.Debug().
			Str("name", r.Name()).
			Str("kind", r.Kind()).
			Stringer("label", label.New("", args.Rel, r.Name())).
			Msg("generated rule")
	}

	for _, r := range args.OtherEmpty {
		dc.Trace().
			Str("name", r.Name()).
			Str("kind", r.Kind()).
			Stringer("label", label.New("", args.Rel, r.Name())).
			Msg("empty rule")
	}

	current := time.Now()
	diff := current.Sub(dl.prev)
	elapsed := current.Sub(dl.start)
	dl.prev = current

	if dc.generaterulesSlowWarnDuration != 0 && diff > dc.generaterulesSlowWarnDuration {
		dc.Warn().
			Str("dir", args.Rel).
			Stringer("t", diff).
			Int("total-rules", len(args.OtherGen)).
			Int("total-files", len(args.RegularFiles)).
			Msgf("slow %s", diff.Round(time.Millisecond))
	}

	dc.Debug().
		Str("label", label.New("", args.Rel, "all").String()).
		Int("rule-count", len(args.OtherGen)).
		Int("file-count", len(args.RegularFiles)).
		Msgf("generated in %s", diff.Round(time.Millisecond))

	if dc.showTotalElapsedTimeMessages {
		dc.Info().
			Str("elapsed", elapsed.Round(time.Millisecond).String()).
			Str("dir", args.Rel).
			Stringer("t", elapsed).
			Msgf("time")
	}

	return language.GenerateResult{}
}
