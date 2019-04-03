#!/bin/bash

cur=$1
new=$2

if [ "$cur" = "" ]; then
    echo -e "Must specify Current tags (for delete): Example command should be ./release.sh 1.0 2.0"
    exit 1
fi

if [ "$new" = "" ]; then
    echo -e "Must specify New tags (for create): Example command should be ./release.sh 1.0 2.0"
    exit 1
fi

git push
git tag -a $new -m "Release $new"
git tag -d $cur
git push origin :$cur
git push origin $new