#!/bin/sh
# Make sure to index data into Elasticsearch first!
#
# Create runtime field for GitHub step duration calculated out of completed_at
# and started_at
# https://www.elastic.co/guide/en/elasticsearch/reference/current/runtime.html
# Adjust if you do want to use anything other then base auth

USER="$1"
PW="$2"
curl -X PUT -u "${USER}:${PW}" 'http://localhost:9200/steps/_mapping/?pretty' \
    -H 'Content-Type: application/json' -d'
{
   "runtime": {
     "duration": {
       "type": "long",
        "script": {
          "source": "emit(doc[\u0027completed_at\u0027].value.millis - doc[\u0027started_at\u0027].value.millis)"
      }
     } 
   } 
}
'
