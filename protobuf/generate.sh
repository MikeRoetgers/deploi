#!/bin/bash

protoc -I ./ -I ../vendor/ server.proto --go_out=plugins=grpc:.
