#!/bin/bash

docker build -t ind .
docker run -d ind
