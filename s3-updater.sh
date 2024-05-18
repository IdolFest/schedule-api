#!/bin/sh

# Set up to be run by cron (eg, every 15 minutes -- */15 * * * *)
# Will pull the schedule from the running schedule api, and upload to the S3 url provided in $1
# Checks for modifications and does not upload if there are none, to reduce s3 requests
#
# Usage: /s3-updater.sh
# Set the following required environment variables in ./s3-updater-conf.sh (gitignored, chmod 700)
# - S3_URL (example: s3://test-bucket/schedule)
# - AWS_ACCESS_KEY_ID
# - AWS_SECRET_ACCESS_KEY

# Need to clear cache? rm -r /tmp/schedule-api-s3

. $(dirname "$0")/s3-updater-conf.sh
mkdir /tmp/schedule-api-s3 2>/dev/null || true
cp /tmp/schedule-api-s3/schedule.json /tmp/schedule-api-s3/schedule-last.json 2>/dev/null || true
/usr/bin/curl -s http://localhost:8080/schedule-by-time > /tmp/schedule-api-s3/schedule.json

if cmp -s "/tmp/schedule-api-s3/schedule.json" "/tmp/schedule-api-s3/schedule-last.json"; then
    # Do nothing, they're equal
    echo "Schedule Identical, not updating"
else
    AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY /usr/bin/aws s3 cp --quiet /tmp/schedule-api-s3/schedule.json $S3_URL/schedule.json
    echo "Schedule updated"
fi
