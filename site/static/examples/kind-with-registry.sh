#!/bin/sh
set -o errexit

# 1. Create registry container unless it already exists
reg_name='kind-registry'
reg_port='5001'
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi

# 2. Set containerd hosts config to remap localhost:${reg_port} on the cluster side
# to the registry container.
#
# This is necessary because localhost resolves to loopback addresses that are
# network-namespace local.
# IE localhost in the container is not localhost on the host.
#
# We want a consistent name that works from both ends, so we tell containerd to
# alias localhost:${reg_port} when pulling images

# 3. Ensure config dir with containerd hosts config
# See: https://github.com/containerd/containerd/blob/main/docs/hosts.md
# TODO: consider built-in kind support for a kind config dir instead
KIND_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/kind/"
HOSTS_DIR="${KIND_DIR}/certs.d"

SERVER_DIR="${HOSTS_DIR}/localhost:${reg_port}"
mkdir -p "${SERVER_DIR}"

cat <<EOF >"${HOSTS_DIR}/localhost/hosts.toml"
server = "http://${reg_name}:5000"
EOF

# 4. Create cluster, with this hosts config enabled
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- extraMounts:
  - hostPath: "${HOSTS_DIR}"
    containerPath: /etc/containerd/certs.d
containerdConfigPatches:
- |-
  # Ensure config_path is enabled so the config below is respected
  # TODO: kind will eventually enable this by default and this patch will
  # be unnecessary.
  #
  # See:
  # https://github.com/kubernetes-sigs/kind/issues/2875
  # https://github.com/containerd/containerd/blob/main/docs/cri/config.md#registry-configuration
  [plugins."io.containerd.grpc.v1.cri".registry]
    config_path = "/etc/containerd/certs.d"
EOF

cat <<EOF | docker exec kind-control-plane cat - /etc/containerd/certs.d/

# 5. Connect the registry to the cluster network if not already connected
# This allows kind to bootstrap the network but ensures they're on the same network
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${reg_name}")" = 'null' ]; then
  docker network connect "kind" "${reg_name}"
fi

# 6. Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF
