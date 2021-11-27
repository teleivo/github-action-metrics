#!/bin/sh
# Make sure to index data into Elasticsearch first!
#
# Create Kibana index patterns for searching
# https://www.elastic.co/guide/en/kibana/master/index-patterns-api-create.html
# Create runtime field for GitHub step duration calculated out of completed_at
# and started_at
# https://www.elastic.co/guide/en/elasticsearch/reference/current/runtime.html
# Adjust if you do want to use anything other then base auth

ELASTIC_USER="$1"
PW="$2"

# Note that runs also have a created_at field. In case you want it to be the
# timeFieldName
curl -i --silent --show-error -X POST -u "${ELASTIC_USER}:${PW}" "localhost:5601/api/index_patterns/index_pattern" \
    -H 'kbn-xsrf: true' -H 'Content-Type: application/json' -d'
{
  "override": true,
  "refresh_fields": true,
  "index_pattern": {
    "id": "runs",
    "title": "runs",
    "timeFieldName": "run_started_at"
  }
}
'
echo "\nCreated index-pattern for runs\n"

curl -i --silent --show-error -X POST -u "${ELASTIC_USER}:${PW}" "localhost:5601/api/index_patterns/index_pattern" \
    -H 'kbn-xsrf: true' -H 'Content-Type: application/json' -d'
{
  "override": true,
  "refresh_fields": true,
  "index_pattern": {
    "id": "jobs",
    "title": "jobs",
    "timeFieldName": "started_at",
    "runtimeFieldMap": {
      "duration": {
        "type": "long",
        "script": {
          "source": "emit(doc[\u0027completed_at\u0027].value.toInstant().toEpochMilli() - doc[\u0027started_at\u0027].value.toInstant().toEpochMilli())"
        }
      }
    },
    "fieldFormats": {
      "duration": {
        "id": "duration",
        "params": {
          "inputFormat": "milliseconds",
          "outputFormat": "asMinutes",
          "useShortSuffix": true,
          "showSuffix": true
        }
      }
    }
  }
}
'
echo "\nCreated index-pattern for jobs\n"

curl -i --silent --show-error -X POST -u "${ELASTIC_USER}:${PW}" "localhost:5601/api/index_patterns/index_pattern" \
    -H 'kbn-xsrf: true' -H 'Content-Type: application/json' -d'
{
  "override": true,
  "refresh_fields": true,
  "index_pattern": {
    "id": "steps",
    "title": "steps",
    "timeFieldName": "started_at",
    "runtimeFieldMap": {
      "duration": {
        "type": "long",
        "script": {
          "source": "emit(doc[\u0027completed_at\u0027].value.toInstant().toEpochMilli() - doc[\u0027started_at\u0027].value.toInstant().toEpochMilli())"
        }
      }
    },
    "fieldFormats": {
      "duration": {
        "id": "duration",
        "params": {
          "inputFormat": "milliseconds",
          "outputFormat": "asMinutes",
          "useShortSuffix": true,
          "showSuffix": true
        }
      }
    }
  }
}
'
echo "\nCreated index-pattern for steps\n"
