#!/bin/bash

# This script make a lot of assumptions and has no error handling


BIN_DIR=`dirname "$0"`
cd $BIN_DIR/../..

BASE_DIR=`pwd`

echo "Base Dir: " $BASE_DIR

rm -rf $BASE_DIR/package
rm -rf $BASE_DIR/OpenDNSMonitor
mkdir $BASE_DIR/package

cp README.md $BASE_DIR/package
cp runit.sh $BASE_DIR/package

###################
# Build and Package CLI for OSX
###################
echo "Building for OSX"
export GOOS="darwin"

mkdir $BASE_DIR/package/$GOOS
go build
mv OpenDNS_Monitor $BASE_DIR/package/$GOOS

###################
# Build and Package CLI for Linux
###################


echo "Building for Linux"
export GOOS="linux"

mkdir $BASE_DIR/package/$GOOS
go build
mv OpenDNS_Monitor $BASE_DIR/package/$GOOS

##################
# Build and Package CLI for Windowss
###################

echo "Building for Windows"
export GOOS="windows"

mkdir $BASE_DIR/package/$GOOS
go build
mv OpenDNS_Monitor.exe $BASE_DIR/package/$GOOS


###################
# Done!!!
###################
mv package OpenDNSMonitor
echo "Done..."
echo ""
echo "See $BASE_DIR/OpenDNSMonitor for binary files"
open $BASE_DIR/OpenDNSMonitor


