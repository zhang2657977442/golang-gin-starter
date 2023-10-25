#!/bin/bash
nohup /export/App/backend --conf=/export/App/config.toml > /export/App/nohup.out 2>&1 &
echo "backend started. pid:[$!]"
echo $! > backend.pid