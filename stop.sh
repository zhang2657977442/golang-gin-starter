#!/bin/sh

if [ ! -e backend.pid ]; then
    echo "ERROR: backend.pid file not exists!!!" 
    exit 1
fi

kill -9 `cat backend.pid`

echo "backend had stopped."
