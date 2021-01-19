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

ARG BASE=kindest/base:v20201130-23777eca
FROM $BASE

# copy in static files
COPY files/ /

ARG KUBERNETES_RELEASE_VERSION="v1.20.2"
RUN echo -n "Writing Kubernetes Version File ..." \
    && echo "$KUBERNETES_RELEASE_VERSION" >/kind/version \
 && echo "Downloading Kubernetes Release Binaries ..." \
    && export ARCH=$(dpkg --print-architecture | sed 's/ppc64el/ppc64le/' | sed 's/armhf/arm/') \
    && export BINARY_BASE_URL="https://dl.k8s.io/${KUBERNETES_RELEASE_VERSION}/bin/linux/${ARCH}" \
    && curl -sSL --retry 5 --output /usr/bin/kubeadm "${BINARY_BASE_URL}/kubeadm" \
    && chmod +x /usr/bin/kubeadm \
    && chown root:root /usr/bin/kubeadm \
    && curl -sSL --retry 5 --output /usr/bin/kubelet "${BINARY_BASE_URL}/kubelet" \
    && chmod +x /usr/bin/kubelet \
    && chown root:root /usr/bin/kubelet \
    && curl -sSL --retry 5 --output /usr/bin/kubectl "${BINARY_BASE_URL}/kubectl" \
    && chmod +x /usr/bin/kubectl \
    && chown root:root /usr/bin/kubectl
