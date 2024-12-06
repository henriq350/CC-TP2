#!/bin/bash


#Enviorment variables
export PATH=$PATH:/user/local/go/bin/go
export GOPATH=$HOME/go
export GOMODCACHE=$GOPATH/pkg/mod
export PATH=$PATH:GOPATH/bin
export HOME=/root
export GOCACHE=/root/.cache/go-build
export TERM=xterm-256color

mkdir -p $GOCACHE