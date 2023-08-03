#!/bin/bash

# These images are for local testing only, if you're looking for CI / deployment scripts, see the `.cloudbuild/`
# directory in this repo.

# For local testing of the slimmed, deployed image
docker build \
  -t api:local \
  --file .cloudbuild/Dockerfile \
  .

# For local testing of the CI process
docker build \
  -t api:local-ci \
  --file .cloudbuild/Dockerfile-Ci \
  .

docker tag api:local-ci gcr.io/fir-sandbox-326008/api:ci
