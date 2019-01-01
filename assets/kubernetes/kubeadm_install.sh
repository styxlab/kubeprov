HOSTNAME=$1
ROLE=$2

echo "Set hostname to $HOSTNAME as $ROLE"
hostnamectl set-hostname $HOSTNAME

IFACE=$(ifconfig -a | grep eth | cut -d ' ' -f 1)
IPV4=$(ip a show $IFACE | grep -m1 inet | tr -s ' ' | cut -d' ' -f3 | cut -d/ -f1)

echo "Make hostname resolvable in /etc/hosts"
echo "$IPV4 $HOSTNAME" >> /etc/hosts

echo "Enable Docker"
systemctl enable docker && systemctl start docker

CNI_VERSION=$(curl -sS https://github.com/containernetworking/plugins/releases/latest | sed 's/[^v0-9.]*//g' | sed  's/^\.\(.*\)\.$/\1/')
echo "Install CNI $CNI_VERSION"
mkdir -p /opt/cni/bin
curl -L "https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-amd64-${CNI_VERSION}.tgz" | tar -C /opt/cni/bin -xz

CRICTL_VERSION=$(curl -sS https://github.com/kubernetes-sigs/cri-tools/releases/latest | sed 's/[^v0-9.]*//g' | sed  's/^\.\(.*\)\.$/\1/')
mkdir -p /opt/bin
curl -L "https://github.com/kubernetes-incubator/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz" | tar -C /opt/bin -xz

RELEASE="$(curl -sSL https://dl.k8s.io/release/stable.txt)"
echo "Install Kubernetes $RELEASE"
mkdir -p /opt/bin
cd /opt/bin
curl -L --remote-name-all https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/amd64/{kubeadm,kubelet,kubectl}
chmod +x {kubeadm,kubelet,kubectl}
cd /home/core

curl -sSL "https://raw.githubusercontent.com/kubernetes/kubernetes/${RELEASE}/build/debs/kubelet.service" | sed "s:/usr/bin:/opt/bin:g" > /etc/systemd/system/kubelet.service
mkdir -p /etc/systemd/system/kubelet.service.d
curl -sSL "https://raw.githubusercontent.com/kubernetes/kubernetes/${RELEASE}/build/debs/10-kubeadm.conf" | sed "s:/usr/bin:/opt/bin:g" > /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

systemctl enable kubelet && systemctl start kubelet

if [ "$2" = "master" ]; then
	/opt/bin/kubeadm init --apiserver-advertise-address=$IPV4  --pod-network-cidr=192.168.0.0/16 --ignore-preflight-errors=NumCPU
	mkdir -p $HOME/.kube
	cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
	chown $(id -u):$(id -g) $HOME/.kube/config

	kubectl apply -f https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/rbac-kdd.yaml
	kubectl apply -f https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml
else
fi
