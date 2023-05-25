#!/bin/bash

docker build \
  -t cors-proxy:local \
  --file .cloudbuild/Dockerfile \
  .
