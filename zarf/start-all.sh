#!/bin/sh
./admin migrate
sleep 5s
./algosearch & ./metrics && node_modules/.bin/next start
#node_modules/.bin/next start

