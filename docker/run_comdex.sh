#!/bin/sh

if test -n "$1"; then
    # need -R not -r to copy hidden files
    cp -R "$1/.comdex" /root
fi

mkdir -p /root/log
comdex start --rpc.laddr tcp://0.0.0.0:26657 --minimum-gas-prices 0.0001ucmdx --trace