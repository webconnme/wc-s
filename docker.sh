#!/bin/sh

docker run --rm -it --volume $(pwd):/work webconn/nc-dev
