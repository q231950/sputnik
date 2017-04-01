go_library(
    name = "go_default_library",
    srcs = ["requestmanager.go"],
    visibility = ["//visibility:public"],
    deps = [
        "//keymanager:go_default_library",
        "@com_github_Sirupsen_logrus//:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["requestmanager_test.go"],
    library = ":go_default_library",
    deps = ["//keymanager/mocks:go_default_library"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
