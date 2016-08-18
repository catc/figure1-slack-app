#!/bin/sh

OUTPUT="fig1-oembed"

if [ -e "./$OUTPUT" ];
then
	echo "Suffixing old build wiht '-old'"
	mv $OUTPUT $OUTPUT-old
fi

go build -o $OUTPUT *.go

echo "Done building as '$OUTPUT'"