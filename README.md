# bazel-gazelle-debug

A [gazelle](https://github.com/bazelbuild/bazel-gazelle) extension that helps
debug what gazelle is doing.

## Usage

```python
load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle", "gazelle_binary")
load("@com_github_bazelbuild_buildtools//buildifier:def.bzl", "buildifier")

# --- show debugging output ---
# gazelle:log_level debug

# --- show summary of total time on .Info ---
# gazelle:show_total_elapsed_time_info_messages true

# --- warn about packages that take more than 1s to generate ---
# gazelle:generaterules_slow_warn_duration 1s

gazelle_binary(
    name = "gazelle-debug",
    languages = DEFAULT_LANGUAGES + ["//language/debug"],
)

gazelle(
    name = "gazelle",
    gazelle = "//:gazelle-debug",
)
```

Here's the output of `bazel run //:gazelle` with `gazelle:log_level debug` and
`gazelle:generaterules_slow_warn_duration 1ms` on this repo:

```
11:50AM DBG configuring directive dir= key=show_total_elapsed_time_info_messages lang=debug value=true
11:50AM DBG configuring directive dir= key=generaterules_slow_warn_duration lang=debug value=1ms
11:50AM DBG visiting dir=language lang=debug
11:50AM DBG visiting dir=language/debug lang=debug
11:50AM DBG read dir dir=language/debug file=BUILD.bazel lang=debug
11:50AM DBG read dir dir=language/debug file=config.go lang=debug
11:50AM DBG read dir dir=language/debug file=debugconfig.go lang=debug
11:50AM DBG read dir dir=language/debug file=fix.go lang=debug
11:50AM DBG read dir dir=language/debug file=generate.go lang=debug
11:50AM DBG read dir dir=language/debug file=kinds.go lang=debug
11:50AM DBG read dir dir=language/debug file=lang.go lang=debug
11:50AM DBG read dir dir=language/debug file=resolve.go lang=debug
11:50AM DBG generated rule kind=go_library label=//language/debug lang=debug name=debug
11:50AM WRN slow 10ms dir=language/debug lang=debug t=9.993983ms total-files=8 total-rules=1
11:50AM DBG generated in 10ms file-count=8 label=//language/debug:all lang=debug rule-count=1
11:50AM INF time dir=language/debug elapsed=10ms lang=debug t=9.994174ms
11:50AM DBG generated in 0s file-count=0 label=//language:all lang=debug rule-count=0
11:50AM INF time dir=language elapsed=10ms lang=debug t=10.132885ms
11:50AM DBG read dir dir= file=.gitignore lang=debug
11:50AM DBG read dir dir= file=BUILD.bazel lang=debug
11:50AM DBG read dir dir= file=README.md lang=debug
11:50AM DBG read dir dir= file=WORKSPACE lang=debug
11:50AM DBG read dir dir= file=bazel-bazel-gazelle-debug lang=debug
11:50AM DBG read dir dir= file=bazel-bin lang=debug
11:50AM DBG read dir dir= file=bazel-out lang=debug
11:50AM DBG read dir dir= file=bazel-testlogs lang=debug
11:50AM DBG read dir dir= file=go.mod lang=debug
11:50AM DBG read dir dir= file=go.sum lang=debug
11:50AM DBG read dir dir= file=go_deps.bzl lang=debug
11:50AM DBG read dir dir= file=workspace_deps.bzl lang=debug
11:50AM WRN slow 7ms dir= lang=debug t=6.528314ms total-files=12 total-rules=0
11:50AM DBG generated in 7ms file-count=12 label=//:all lang=debug rule-count=0
11:50AM INF time dir= elapsed=17ms lang=debug t=16.661199ms
11:50AM DBG resolving label that provides 'import' impLang=go import=github.com/bazelbuild/bazel-gazelle/config lang=debug resolve=go
11:50AM DBG resolving label that provides 'import' impLang=go import=github.com/bazelbuild/bazel-gazelle/label lang=debug resolve=go
11:50AM DBG resolving label that provides 'import' impLang=go import=github.com/bazelbuild/bazel-gazelle/language lang=debug resolve=go
11:50AM DBG resolving label that provides 'import' impLang=go import=github.com/bazelbuild/bazel-gazelle/repo lang=debug resolve=go
11:50AM DBG resolving label that provides 'import' impLang=go import=github.com/bazelbuild/bazel-gazelle/resolve lang=debug resolve=go
11:50AM DBG resolving label that provides 'import' impLang=go import=github.com/bazelbuild/bazel-gazelle/rule lang=debug resolve=go
11:50AM DBG resolving label that provides 'import' impLang=go import=github.com/rs/zerolog lang=debug resolve=go
```


## Dependencies

Other than gazelle itself, this package requires [github.com/rs/zerolog](https://github.com/rs/zerolog).