#!/bin/bash

set -e
cd "$( dirname "${BASH_SOURCE[0]}" )/../"

git describe --tags --abbrev=0 | sed s/^v//
