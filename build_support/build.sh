#!/bin/bash

# has to be run from 'build_support' subfolder
cd ..
echo "compiling BlackJack_cli and BlackJack_service ..."
env GOOS=linux GOARCH=arm GOARM=6 go build -o build/BlackJack_service cmd/BlackJack_service/BlackJack_service.go
env GOOS=linux GOARCH=arm GOARM=6 go build -o build/BlackJack_cli cmd/BlackJack_cli/BlackJack_cli.go

echo "compiling web client to JavaScript ..."
cd web_client
gopherjs build -o ../build/webapp.js

echo "...Results stored in ./build directory"
echo
echo "On BlackJack ALOA the compiled files have to be placed at the following"
echo "locations:"
echo
echo "    /usr/local/bin/BlackJack_cli"
echo "    /usr/local/bin/BlackJack_service"
echo "    /usr/local/BlackJack/www/webapp.js"
echo "    /usr/local/BlackJack/www/webapp.js.map"

