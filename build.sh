#!/bin/bash
echo "dkr_01 build script start"
# move to /plugin dir
cd plugin || exit 1

# check for plugin-build's existence and make it read-writable if it exists
if [ -d ./plugin-build ]
then
  chmod -R +rw ./plugin-build
fi
# run make all
make all
echo "dkr_01 build script done"