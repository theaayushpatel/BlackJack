#!/bin/bash

# dependencies for the web app
gopherjs build -o ../build/webapp.js #main.go
#scp ../build/webapp* pi@raspberrypi.local:/usr/local/BlackJack/www/
#scp ../dist/www/index.html pi@raspberrypi.local:/usr/local/BlackJack/www/
#scp ../dist/www/BlackJack.css pi@raspberrypi.local:/usr/local/BlackJack/www
scp ../build/webapp* 172.16.0.1:/usr/local/BlackJack/www/
scp ../dist/www/index.html 172.16.0.1:/usr/local/BlackJack/www/
scp ../dist/www/BlackJack.css 172.16.0.1:/usr/local/BlackJack/www/
