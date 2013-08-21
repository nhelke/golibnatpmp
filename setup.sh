#!/bin/bash

# Setup script highly inspired by
# http://areyoufuckingcoding.me/2013/04/15/dont-repeat-yourself-with-setting-up/

# Project specific setup
VERSION="libnatpmp-20120821"
FILES="declspec.h {getgateway,natpmp}.{c,h}"
UPSTREAMURL="http://miniupnp.free.fr/files/download.php?file=$VERSION.tar.gz"

# Enable some bash magic that pipe commands will fail.
# $ info bash > Builtin Commands > The Set builtin
set -e
set -o pipefail

function prefixed {
    sed -e "s/^/       /"
}

function assert_command {
    cmd="$1"

    command -v $cmd >/dev/null 2>&1 || {
        echo "$cmd: Command required but not found. Aborting." >&2;
        exit 1;
    }

    echo "$cmd: $(command -v $cmd)";
}

function assert_commands {
    commands="$@"

    for cmd in $commands; do
        assert_command $cmd
    done
}

# Check that we are running in the script's directory, i.e. package root
function assert_cwd {
	DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
	if [ $DIR == pwd ]
	then
		echo "This script should be run from the package base directory: $DIR" >&2;
		echo "Aborting." >&2;
		exit 2
	fi
}

function install_dependencies {
    # From http://unix.stackexchange.com/a/84980
    TMPFILE=$(mktemp -d 2>/dev/null || mktemp -dt 'libnatpmp.tgz')
    pushd $TMPFILE > /dev/null
    curl $UPSTREAMURL -sLo $VERSION.tgz
    tar -xf $VERSION.tgz
    pushd $VERSION > /dev/null
    eval cp $FILES $(dirs +2)
    popd > /dev/null
    popd > /dev/null
    # go get || exit $?
}

assert_commands sed mktemp 2>&1 > /dev/null

echo "-----> Checking for required commands and tools"
assert_commands go curl tar 2>&1 | prefixed

echo "-----> Checking CWD"
assert_cwd 2>&1 | prefixed

# echo "-----> Installing dependencies"
echo "-----> Fetching libnatpmp C source and extracting required source files only"
install_dependencies 2>&1 | prefixed

echo "-----> All set!"
