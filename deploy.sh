#!/bin/sh

git fetch origin
git reset --hard origin/master
go build
/etc/init.d/canihave restart
