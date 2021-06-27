#!/usr/bin/env bash

export URL="https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=$CANIHAVE_MAXMIND_KEY&suffix=tar.gz"

echo $URL
curl --output GeoLite2-Country.tar.gz $URL | tar -xz
