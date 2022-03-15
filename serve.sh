#!/bin/sh

npm start --prefix ./frontend & cd .
go run main.go &

wait