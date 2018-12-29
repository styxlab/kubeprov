#!/bin/bash

while ! nmap -Pn -p 22 $1 |grep open &>/dev/null; do sleep 2; done