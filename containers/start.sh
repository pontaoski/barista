#!/bin/sh
mkdir -p /var/lib/dbus
dbus-uuidgen > /var/lib/dbus/machine-id

export $(dbus-launch)

cd QueryKit-*
./QueryKit.py &

cd ..
mkdir -p storage
./barista