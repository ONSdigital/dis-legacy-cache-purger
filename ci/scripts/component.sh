#!/bin/bash -eux

pushd dis-legacy-cache-purger
  make test-component
popd
