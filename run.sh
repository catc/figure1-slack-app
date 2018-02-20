#! /usr/bin/env bash

IMAGE_NAME="fig1slack"
SERVICE_NAME="fig1-slack"
PORT=3400

function build {
	docker build -t ${IMAGE_NAME} .
}

function run {
	echo "Starting up figure 1 slackbot service!"
	docker service create -p ${PORT}:${PORT} --name ${SERVICE_NAME} ${IMAGE_NAME}
}

if [ $# -eq 0 ]; then
	echo "must specify 'build' or 'run'"
elif [ "$1" == "build" ]; then
	build
elif [ "$1" == "run" ]; then
	run
else
	exit_with_err
fi