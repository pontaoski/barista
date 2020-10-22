#!/bin/sh
mkdir -p /var/lib/dbus
dbus-uuidgen > /var/lib/dbus/machine-id

export $(dbus-launch)

cd /barista/QueryKit-*/
./QueryKit.py &

cd /barista
mkdir -p storage
./barista