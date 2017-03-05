#!/bin/bash

APPNAME=dcenv

apt-get update
apt-get install -y zip
go get github.com/aktau/github-release

CUR=`readlink -f $(dirname $0)`
pushd ${CUR}
mkdir tmp

#compile
#VERSION="1.1.1"
#VERSION="$(git describe --tags $1)"
VERSION="$(git describe --abbrev=0)"
github-release release -u nak1114 -r dcenv -t ${VERSION} --name ${VERSION} 

for OS in "freebsd" "linux" "darwin" "windows"; do
  for ARCH in "386" "amd64"; do

    EXEC=${CUR}/tmp/${APPNAME}/files/${APPNAME}
    [ ${OS} = "windows" ] && EXEC=${EXEC}.exe
    
    cp -r ${CUR}/item ${CUR}/tmp/${APPNAME}
    
    
    pushd ..
    GOOS=${OS} CGO_ENABLED=0 GOARCH=${ARCH} go build -ldflags "-X main.Version=$VERSION" -o ${EXEC}
    popd

    #compress
    pushd tmp/${APPNAME}
    if [ ${OS} = "windows" ] ; then
      ARCHIVE="${APPNAME}-${VERSION}-${OS}-${ARCH}.zip"
      zip -r ../${ARCHIVE} *
    else
      ARCHIVE="${APPNAME}-${VERSION}-${OS}-${ARCH}.tar.gz"
      tar -czf ../${ARCHIVE} *
    fi
    popd
    rm -rf ${CUR}/tmp/${APPNAME}

    #release
    echo ${ARCHIVE}
    github-release upload -u nak1114 -r dcenv -t ${VERSION} -f ${CUR}/tmp/${ARCHIVE} -n ${ARCHIVE}
  done
done
popd
#for OS in "linux" "windows"; do
#  for ARCH in "amd64"; do
#for OS in "freebsd" "linux" "darwin" "windows"; do
#  for ARCH in "386" "amd64"; do
# powershell Expand-Archive dcenv-1.1.1-windows-amd64.zip  -DestinationPath %DCENV_HOME%
# tar xvfz dcenv-1.1.1-linux-amd64.tar.gz -C $DCENV_HOME
