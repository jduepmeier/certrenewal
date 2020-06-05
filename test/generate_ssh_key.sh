#!/bin/bash
DIR="${1:-bin/test}"
mkdir -p $DIR
ssh-keygen -t ed25519 -f "$DIR/ssh_host_ed25519_key" -q -N ""