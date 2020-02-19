#!/bin/bash

curl -X POST -s -H "Content-type: application/json" -d '{"name":"shane"}' localhost:8080/rooms/create

curl -X GET -s -H "Content-type: application/json" localhost:8080/room/$ROOM_CODE

curl -X GET -s -H "Content-type: application/json" localhost:8080/rooms

curl -X POST -s -H "Content-type: application/json" -d '{"name": "chelsea"}' localhost:8080/room/$ROOM_CODE/join

curl -X POST -s -H "Content-type: application/json" -d '{"name": "chelsea"}' localhost:8080/room/$ROOM_CODE/quit