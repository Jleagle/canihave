#!/bin/sh

echo "### Pulling"
git fetch origin
git reset --hard origin/master

echo "### Building"
dep ensure
go build

echo "### Rollbar"
curl https://api.rollbar.com/api/1/deploy/ \
  -F access_token=${CANIHAVE_ROLLBAR_PRIVATE} \
  -F environment=${ENV} \
  -F revision=$(git log -n 1 --pretty=format:"%H") \
  -F local_username=Jleagle \
  --silent > /dev/null

echo "### Restarting"
/etc/init.d/canihave restart
