#!/bin/sh

git fetch origin
git reset --hard origin/master

godep restore
go build

/etc/init.d/canihave restart
