#!/bin/bash

sudo cp contrib/systemd/matchbox-on-coreos.service /etc/systemd/system/matchbox.service

sudo systemctl daemon-reload
sudo systemctl start matchbox
sudo systemctl enable matchbox

systemctl status matchbox