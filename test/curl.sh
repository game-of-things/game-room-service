#!/bin/bash

curl -X POST -H "Content-type: application/json" -d '{"name":"shane"}' localhost:8080/rooms/create

curl -X GET -H "Content-type: application/json" localhost:8080/room/ZSCV

curl -X GET -H "Content-type: application/json" localhost:8080/rooms

curl -X POST -H "Content-type: application/json" -d '{"name": "chelsea"}' localhost:8080/room/KPFD/join