go_prefix("github.com/q231950/sputnik")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    visibility = ["//visibility:private"],
    deps = ["//cmd:go_default_library"],
)

go_binary(
    name = "sputnik",
    library = ":go_default_library",
    visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_prefix")
