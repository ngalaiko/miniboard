workspace(
    name = "miniboard",
    managed_directories = {
        "@npm": ["web/node_modules"],
    },
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Golang rules.
http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/0.19.3/rules_go-0.19.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/0.19.3/rules_go-0.19.3.tar.gz",
    ],
    sha256 = "313f2c7a23fecc33023563f082f381a32b9b7254f727a7dd2d6380ccc6dfe09b",
)
load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

# Gazelle.
http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/0.18.1/bazel-gazelle-0.18.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.18.1/bazel-gazelle-0.18.1.tar.gz",
    ],
    sha256 = "be9296bfd64882e3c08e3283c58fcb461fa6dd3c171764fcc4cf322f60615a9b",
)
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
gazelle_dependencies()

# Protobuf compiler.
http_archive(
    name = "com_google_protobuf",
    urls = [
        "https://github.com/protocolbuffers/protobuf/archive/v3.9.0.zip",
    ],
    strip_prefix = "protobuf-3.9.0",
    sha256 = "8eb5ca331ab8ca0da2baea7fc0607d86c46c80845deca57109a5d637ccb93bb4",
)
load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")
protobuf_deps()

# Docker image.
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "e513c0ac6534810eb7a14bf025a0f159726753f97f74ab7863c650d26e01d677",
    strip_prefix = "rules_docker-0.9.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.9.0/rules_docker-v0.9.0.tar.gz"],
)
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)
container_repositories()
load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)
_go_image_repos()
load("@io_bazel_rules_docker//toolchains/docker:toolchain.bzl",
    docker_toolchain_configure="toolchain_configure"
)

# web dependencies
http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "ecaa54955b314b5e33948bd8f39e35c35ee89e905d8de1c03868100293510573",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/1.2.1/rules_nodejs-1.2.1.tar.gz"],
)
load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")
yarn_install(
    name = "npm",
    package_json = "//web:package.json",
    yarn_lock = "//web:yarn.lock",
)
load("@npm//:install_bazel_dependencies.bzl", "install_bazel_dependencies")
install_bazel_dependencies()

# gRPC-Web
http_archive(
    name = "rules_proto_grpc",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/1.0.1.tar.gz"],
    sha256 = "497225bb586e8f587e139c55b0f015e93bdddfd81902985ce24623528dbe31ab",
    strip_prefix = "rules_proto_grpc-1.0.1",
)
load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_toolchains", "rules_proto_grpc_repos")
rules_proto_grpc_toolchains()
rules_proto_grpc_repos()
# commonjs_grpc_compile
load("@rules_proto_grpc//github.com/grpc/grpc-web:repositories.bzl", rules_proto_grpc_grpc_web_repos="grpc_web_repos")
rules_proto_grpc_grpc_web_repos()
load("@io_bazel_rules_closure//closure:repositories.bzl", "rules_closure_dependencies", "rules_closure_toolchains")
rules_closure_dependencies(
    omit_com_google_protobuf = True,
)
rules_closure_toolchains()

# Golang app's deps.
load("//server:repositories.bzl", "go_repositories")
go_repositories()
