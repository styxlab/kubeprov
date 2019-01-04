#!/bin/bash

HOSTNAME=$1

echo "Set hostname to $HOSTNAME as $ROLE"
hostnamectl set-hostname $HOSTNAME

IFACE=$(ifconfig -a | grep eth | cut -d ' ' -f 1)
IPV4=$(ip a show $IFACE | grep -m1 inet | tr -s ' ' | cut -d' ' -f3 | cut -d/ -f1)

echo "Make hostname resolvable in /etc/hosts"
echo "$IPV4 $HOSTNAME" >> /etc/hosts

echo "Enable Docker"
systemctl enable docker && systemctl start docker
