workspace(name = "build_stack_bazel_gazelle_debug")

load("//:workspace_deps.bzl", "workspace_deps")

workspace_deps()

# -----------------------------------------
# rules_go
# -----------------------------------------

load("@io_bazel_rules_go//go:deps.bzl", "go_download_sdk", "go_register_toolchains", "go_rules_dependencies")

go_download_sdk(
    name = "go_sdk_linux",
    goarch = "amd64",
    goos = "linux",
    version = "1.16.2",
)

go_download_sdk(
    name = "go_sdk_darwin",
    goarch = "amd64",
    goos = "darwin",
    version = "1.16.2",
)

go_register_toolchains()

go_rules_dependencies()

# -----------------------------------------
# gazelle
# -----------------------------------------

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

# gazelle:repo bazel_gazelle

load("//:go_deps.bzl", "fetch_go_deps")

# gazelle:repository_macro go_deps.bzl%fetch_go_deps
fetch_go_deps()

gazelle_dependencies()

# -----------------------------------------
# buildifier
# -----------------------------------------

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
