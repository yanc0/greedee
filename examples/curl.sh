#!/bin/sh

set -e
set -u

### CHANGE ME ###
EVENTS_URL="http://127.0.0.1:9223/events"
USER="user"
PASSWORD="pass"
CURL=$(which curl)

name="toto"
ttl=0
status=0
source="curl"
description="blah blah"
#################

$CURL -d "{ \"name\": \"$name\", \"ttl\": $ttl,
            \"status\": $status, \"source\": \"$source\",
            \"description\": \"$description\"}" \
            $EVENTS_URL -u "$USER:$PASSWORD"   
