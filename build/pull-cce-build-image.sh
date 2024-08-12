#! /bin/bash
# Copyright 2024 Baidu, Inc.
set -ex


BUILD_IMAGES=(
    "build-image/debian-iptables:buster-v1.6.7"
    "build-image/kube-cross:v1.15.15-legacy-1"
    "build-image/go-runner:v2.3.1-go1.15.15-buster.0"
)

GCR_REGISTRY=k8s.gcr.io
CCE_REGISTRY=registry.baidubce.com/cce-plugin-dev/kubernetes

for image in ${BUILD_IMAGES[@]}; do
    docker pull $GCR_REGISTRY/$image && \
    docker tag $GCR_REGISTRY/$image $CCE_REGISTRY/$image && \
    docker push $CCE_REGISTRY/$image
done