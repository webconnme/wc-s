#!/bin/sh

go get github.com/go-martini/martini
go get github.com/martini-contrib/binding
go get github.com/martini-contrib/render
go get github.com/martini-contrib/cors

go get --tags zmq_4_x github.com/alecthomas/gozmq
go get github.com/mikepb/go-serial


cd src
for MODULE in wc-s/*
do
	go install $MODULE
done
