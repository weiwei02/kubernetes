# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This file create the kube-apiserver image.
# we use the hosts platform to apply the capabilities to avoid the need
# to setup qemu for the builder.
FROM registry.baidubce.com/cce-plugin-dev/kubernetes/kube-apiserver-amd64:v1.24.17-baidu-0716-dirty

COPY bin/linux/amd64/kube-apiserver  /usr/local/bin/kube-apiserver
