#!/bin/sh

brew services restart memcached

git fetch origin
git reset --hard origin/master
go build
/etc/init.d/canihave restart
