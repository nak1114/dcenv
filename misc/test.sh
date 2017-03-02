#!/bin/bash
# make test dcenv environment

rm -rf /test
mkdir -p /test/path/to/cur/child/path
mkdir -p /test/path/not/to/cur

CUR=`dirname $0`

cp -r ${CUR}/item /test
cp -r ${CUR}/item /test/path

cp ${CUR}/test/.dcenv_linux_global   /test/item/files/.dcenv_bash
cp ${CUR}/test/.dcenv_linux_localcur /test/path/to/cur/.dcenv_bash
cp ${CUR}/test/.dcenv_linux_localto  /test/path/to/.dcenv_bash
