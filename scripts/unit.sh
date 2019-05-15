#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

# Packages to exclude
PKGS=`go list github.com/trustbloc/sidetree-core-go/pkg/... 2> /dev/null | \
                                                   grep -v /mocks | \
                                                   grep -v /api/`
echo "Running pkg unit tests..."
go test -count=1 -cover $PKGS -p 1 -timeout=10m
