#!/bin/sh

OUTPUT="figure1-slack-app"

if [ -e "./$OUTPUT" ];
then
	echo "Suffixing old build with '-old'"
	mv $OUTPUT $OUTPUT-old
fi

go build -o $OUTPUT *.go

echo "Done building as '$OUTPUT'"