#!/usr/bin/bash

# usage: setup {node_index}
function setup() {
  sudo apt-get update
  cd /users/patrickm || exit
  wget -P /tmp https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
  sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf /tmp/go1.24.2.linux-amd64.tar.gz
  echo "export node_index=$1" >> /users/patrickm/.bashrc
  echo "export PATH=$PATH:/usr/local/go/bin" >> /users/patrickm/.bashrc
  source /users/patrickm/.bashrc
}