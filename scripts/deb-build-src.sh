#!/bin/bash

set -e
cd "$( dirname "${BASH_SOURCE[0]}" )"

mkdir src
cd src
tar xf ../mainline-kernel-tool_*.orig.tar.gz
dpkg-buildpackage -uc -us -S
