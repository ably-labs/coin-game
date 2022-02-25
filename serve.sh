#!/bin/sh

npm start --prefix ./frontend &
gin --port 8000 --path . --build ./src/server/ --i --all &

wait