# 构建测试镜像，并推送到 CCE dev 镜像仓库
# 注意：这个脚本仅用于快速构建测试包，不可直接用于生产发布
IMAGE_VERSION="20240719a"

export KUBE_DOCKER_REGISTRY=registry.baidubce.com/cce-plugin-dev/kubernetes
export KUBE_BASE_IMAGE_REGISTRY=registry.baidubce.com/cce-plugin-dev/kubernetes/build-image
export KUBE_RELEASE_RUN_TESTS=n

# 使用案例
# 构建并发布镜像
# ./build/cce_build.sh
KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

KUBE_SERVER_IMAGE_TARGETS=cmd/kube-apiserver \
KUBE_SERVER_PLATFORMS="linux/amd64" \
${KUBE_ROOT}/build/run.sh make all WHAT=cmd/kube-apiserver

docker build -t registry.baidubce.com/cce-plugin-dev/kubernetes/kube-apiserver-amd64:${IMAGE_VERSION} -f build/server-image/kube-apiserver/test.Dockerfile ./_output/dockerized
docker push registry.baidubce.com/cce-plugin-dev/kubernetes/kube-apiserver-amd64:${IMAGE_VERSION}