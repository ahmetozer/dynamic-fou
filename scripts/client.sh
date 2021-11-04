#!/usr/bin/env bash
case "${REMOTE_ADDR}:${REMOTE_PORT}" in

"203.0.113.80:65200")
    ifconfig $INTERFACE 192.168.2.2 netmask 255.255.255.0
    iptables -A FORWARD -i  $INTERFACE -j ACCEPT
    iptables -A FORWARD -o  $INTERFACE -j ACCEPT
    ;;

"203.0.113.15:65205")
    ifconfig $INTERFACE 192.168.3.2 netmask 255.255.255.0
    iptables -A FORWARD -i  $INTERFACE -j ACCEPT
    iptables -A FORWARD -o  $INTERFACE -j ACCEPT
    iptables -t nat -A POSTROUTING -o $INTERFACE -j MASQUERADE 
    ;;

*)
    echo >&2 -n "unknown server"
    exit 1
    ;;
esac
