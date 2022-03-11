#!/bin/sh

npm start --prefix ./frontend &
go run main.go &

wait