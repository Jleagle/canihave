#!/bin/sh

git fetch origin
git reset --hard origin/master

dep ensure
go build

/etc/init.d/canihave restart
