export KUBE_DOCKER_REGISTRY=registry.baidubce.com/cce-plugin-dev/kubernetes
export KUBE_BASE_IMAGE_REGISTRY=registry.baidubce.com/cce-plugin-dev/kubernetes/build-image
export KUBE_RELEASE_RUN_TESTS=n

# 使用案例
# 构建并发布镜像
# ./build/cce_build.sh
KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

KUBE_BUILD_TARGETS=cmd/kube-apiserver \
KUBE_BUILD_PLATFORMS="linux/amd64 linux/arm64" \
KUBE_BUILD_CONFORMANCE=n \
${KUBE_ROOT}/build/run.sh "$@"