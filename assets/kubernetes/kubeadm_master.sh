#!/bin/bash

echo "I'm the master"

IFACE=$(ifconfig -a | grep eth | cut -d ' ' -f 1)
IPV4=$(ip a show $IFACE | grep -m1 inet | tr -s ' ' | cut -d' ' -f3 | cut -d/ -f1)

echo $IPV4

sudo /opt/bin/kubeadm init --apiserver-advertise-address=$IPV4  --pod-network-cidr=192.168.0.0/16 --ignore-preflight-errors=NumCPU

echo $HOME
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown core:core $HOME/.kube/config

until $(ncat -z ${IPV4} 6443); do
  echo "Waiting for API server to respond"
  sleep 5
done

echo "kubectl1"
/opt/bin/kubectl get nodes

echo "kubectl2"
/opt/bin/kubectl apply -f https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/rbac-kdd.yaml
/opt/bin/kubectl apply -f https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml

echo "Finished"