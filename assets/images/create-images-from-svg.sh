#!/bin/sh
#
# Script to generate all the .png images from all the .svg images.

THIS_DIR="$( cd "$(dirname "$0")"; pwd -P )"

if [ ! -d SVG ]
then
    echo "Please run this script inside the assets/images directory only."
    exit 1
fi

ls -1 SVG/*.svg | while read SVG
do
    PNG=$(echo ${SVG} | sed 's/\.svg/.png/' | sed 's+SVG/++' )
    for SIZE in 16 32 64 128
    do
        mkdir -p ${SIZE}
	echo "Converting ${SVG} to size ${SIZE}x${SIZE} ..."
        inkscape \
            --without-gui \
            --file ${SVG} \
            --export-png ${SIZE}/${PNG} \
            --export-width ${SIZE} \
            --export-height ${SIZE}
    done
done
