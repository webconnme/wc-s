#!/bin/sh

cd src
for MODULE in wc-s/app/*
do
	echo build $MODULE
	go install $MODULE
done

for MODULE in wc-s/env/*
do
	echo build $MODULE
	go install $MODULE
done
cd ..

cd native
make -C led
cd ..
