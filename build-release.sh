#!/bin/bash
NAME=actions-usage
PLATFORMS="linux darwin windows"
VERSION=0.1.0

go get
for i in $PLATFORMS; do
  rm -rf target/$i; mkdir -p target/$i/$NAME
  GOOS=$i go build -o target/$i/$NAME
  (cd target/$i; tar czvf $NAME-$VERSION-$i.tar.gz $NAME; zip $NAME-$VERSION-$i.zip -r $NAME)
done;
