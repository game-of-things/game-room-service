#!/bin/bash

#HOST_NAME=http://localhost:8080
HOST_NAME=http://game-room-gameofthings.apps-crc.testing:80

curl -X POST -s -H "Content-type: application/json" -d '{"name":"shane"}' $HOST_NAME/rooms/create

curl -X GET -s -H "Content-type: application/json" $HOST_NAME/room/$ROOM_CODE

curl -X GET -s -H "Content-type: application/json" $HOST_NAME/rooms

curl -X POST -s -H "Content-type: application/json" -d '{"name": "matt"}' $HOST_NAME/room/$ROOM_CODE/join

curl -X POST -s -H "Content-type: application/json" -d '{"name": "matt"}' $HOST_NAME/room/$ROOM_CODE/quit