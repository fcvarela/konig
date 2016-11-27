#!/bin/bash

files=`ls *.proto`

protoc -I. \
       -I/usr/local/include \
       --go_out=plugins=grpc:. \
        $files
