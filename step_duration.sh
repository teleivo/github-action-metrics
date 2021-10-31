#!/bin/sh
curl -X PUT 'http://localhost:9200/steps/_mapping/?pretty' \
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

# curl -X PUT "localhost:9200/seats/_mapping?pretty" -H 'Content-Type: application/json' -d'
# {
#   "runtime": {
#     "day_of_week": {
#       "type": "keyword",
#       "script": {
#         "source": "emit(doc[\u0027datetime\u0027].value.getDayOfWeekEnum().toString())"
#       }
#     }
#   }
# }
# '

