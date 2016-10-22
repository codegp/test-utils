#!/bin/bash

set -e

PROJECT_ID=$(gcloud config list project --format "value(core.project)" 2> /dev/null)
pushd $GOPATH/src/github.com/codegp/test-utils/setuprunner
CGO_ENABLED=0 go build
docker build -t gcr.io/$PROJECT_ID/setuprunner .
rm setuprunner
popd
