#!/bin/bash

echo "mount /dev/sda9"
mount /dev/sda9 /mnt

CNI_VERSION=$(curl -sS https://github.com/containernetworking/plugins/releases/latest | sed 's/[^v0-9.]*//g' | sed  's/^\.\(.*\)\.$/\1/')
echo "Install CNI $CNI_VERSION"
mkdir -p /mnt/opt/cni/bin
curl -L "https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-amd64-${CNI_VERSION}.tgz" | tar -C /mnt/opt/cni/bin -xz

CRICTL_VERSION=$(curl -sS https://github.com/kubernetes-sigs/cri-tools/releases/latest | sed 's/[^v0-9.]*//g' | sed  's/^\.\(.*\)\.$/\1/')
mkdir -p /mnt/opt/bin
curl -L "https://github.com/kubernetes-incubator/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz" | tar -C /mnt/opt/bin -xz

RELEASE="$(curl -sSL https://dl.k8s.io/release/stable.txt)"
echo "Install Kubernetes $RELEASE"
mkdir -p /mnt/opt/bin
cd /mnt/opt/bin
curl -L --remote-name-all https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/amd64/{kubeadm,kubelet,kubectl}
chmod +x {kubeadm,kubelet,kubectl}

echo "Get kubelet"
curl -sSL "https://raw.githubusercontent.com/kubernetes/kubernetes/${RELEASE}/build/debs/kubelet.service" | sed "s:/usr/bin:/opt/bin:g" > /mnt/etc/systemd/system/kubelet.service
mkdir -p /mnt/etc/systemd/system/kubelet.service.d
curl -sSL "https://raw.githubusercontent.com/kubernetes/kubernetes/${RELEASE}/build/debs/10-kubeadm.conf" | sed "s:/usr/bin:/opt/bin:g" > /mnt/etc/systemd/system/kubelet.service.d/10-kubeadm.conf

echo "Finished"
