go_library(
    name = "go_default_library",
    srcs = [
        "create.go",
        "identity.go",
        "identitydelete.go",
        "keyid.go",
        "ping.go",
        "print.go",
        "root.go",
        "store.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//keymanager:go_default_library",
        "//requesthandling:go_default_library",
        "@com_github_Sirupsen_logrus//:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")
