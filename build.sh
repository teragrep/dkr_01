#!/bin/bash
echo "dkr_01 build script start"
# move to /plugin dir
cd plugin || exit 1

# check for plugin-build existence and make it rwx if it exist
if [ -d ./plugin-build ]
then
  chmod -R +rwx ./plugin-build
fi
# run make all
make all
echo "dkr_01 build script done"