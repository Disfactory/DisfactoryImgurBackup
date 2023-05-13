#!/bin/bash

pushd backup
for file in *.gif
do
    newname=$(echo $file | cut -d'_' -f3)
    echo "mv $file $newname"
    mv $file $newname
done
popd
