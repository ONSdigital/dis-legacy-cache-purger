#!/bin/bash -eux

pushd dis-legacy-cache-purger
  make build
  cp build/dis-legacy-cache-purger Dockerfile.concourse ../build
popd
