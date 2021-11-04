#!/usr/bin/env bash

case $CLIENT_NAME in

client1)
    ifconfig $INTERFACE 192.168.2.1 netmask 255.255.255.0
    ip ro add 10.0.5.0/24 via 192.168.2.2
    ;;

client2)
    ifconfig $INTERFACE 10.0.6.1 netmask 255.255.255.0
    ip ro add 10.0.7.0/24 via 10.0.6.1
    ;;

*)
    echo >&2 -n "unknown client"
    exit 1
    ;;
esac
