#!/bin/bash
# make local dcenv environment

[ -z "$DCENV_OLDPATH" ] && export DCENV_OLDPATH=${PATH}
[ -z "$DCENV_OLDHOME" ] && export DCENV_OLDHOME=${DCENV_HOME}
cp -r ./misc/item ./.dcenv_home
export DCENV_HOME=$PWD/.dcenv_home
export PATH=${DCENV_HOME}/bin:${DCENV_HOME}/shims:${DCENV_OLDPATH}
if [ -x "`which docker`" ] ; then
  docker run --rm -it -v $(pwd)/misc/go:/go -v $(pwd)/misc/tmp:/tmp -v $(pwd):/go/src/github.com/nak1114/dcenv -w /go/src/github.com/nak1114/dcenv -e GOPATH=/go golang go build -o ${DCENV_HOME}/files/dcenv
else
  go build -o ${DCENV_HOME}/files/dcenv
fi
