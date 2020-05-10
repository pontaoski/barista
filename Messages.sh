#!/bin/sh

xgettext -kI18n -kI18nc:1c,2 `find . -name \*.go` -o messages/barista.pot -j
for pofile in messages/*.po; do
    msgmerge -U $pofile messages/barista.pot
done