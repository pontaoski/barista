#!/bin/sh
mkdir -p /var/lib/dbus
dbus-uuidgen > /var/lib/dbus/machine-id

export $(dbus-launch)

cd /barista/QueryKit-0.1/
./QueryKit.py &

cd /barista
mkdir -p storage
./barista