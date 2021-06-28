#!/usr/bin/env bash

export URL="https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=$CANIHAVE_MAXMIND_KEY&suffix=tar.gz"
export TAR="GeoLite2-Country.tar.gz"
export FILE="GeoLite2-Country.mmdb"

rm ./pkg/location/$FILE
curl --silent --output $TAR $URL
tar --extract -z -v --strip-components 1 --file=$TAR $(tar -t -f $TAR | grep mmdb)
mv ./$FILE ./pkg/location/$FILE
rm $TAR
