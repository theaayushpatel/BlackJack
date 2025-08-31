SHELL := /bin/bash
PATH := /usr/local/go/bin:$(PATH)

all: compile

test:
#	export PATH="$$PATH:/usr/local/go/bin" # put into ~/.profile
	echo $(CURDIR)
	echo $(HOME)

# make dep runs without sudo
dep:
	# sudo apt-get -y install git screen hostapd autossh bluez bluez-tools bridge-utils policykit-1 genisoimage iodine haveged
	# sudo apt-get -y install tcpdump
	# sudo apt-get -y install python-pip python-dev

	# before installing dnsmasq, the nameserver from /etc/resolv.conf should be saved
	# to restore after install (gets overwritten by dnsmasq package)
	# cp /etc/resolv.conf /tmp/backup_resolv.conf
	# sudo apt-get -y install dnsmasq
	# sudo /bin/bash -c 'cat /tmp/backup_resolv.conf > /etc/resolv.conf'

	# python dependencies for HIDbackdoor
	# sudo pip install pycrypto # already present on stretch
	# sudo pip install pydispatcher

	# install go
	# wget https://storage.googleapis.com/golang/go1.10.linux-armv6l.tar.gz
	# sudo tar -C /usr/local -xzf go1.10.linux-armv6l.tar.gz

	export PATH="$$PATH:/usr/local/go/bin"

	# put into ~/.profile
	# ToDo: check if already present
	# echo "export PATH=\$$PATH:/usr/local/go/bin" >> ~/.profile
	# sudo bash -c 'echo export PATH=\$$PATH:/usr/local/go/bin >> ~/.profile'

	# install gopherjs
	go install github.com/gopherjs/gopherjs

	# we don't need protoc + protoc-grpc-web, because the proto file is shipped pre-compiled

	# go dependencies for webapp (without my own)
	#go get google.golang.org/grpc
	#go get -u github.com/improbable-eng/grpc-web/go/grpcweb
	#go get -u github.com/gorilla/websocket

# This target probably needs to be run at least once to get the dependencies on
# the go path. But after that, you probably actually want to run:
# $ cd build_support && ./build.sh && cd ..
# instead, to build with the right GOOS and GOARCH settings.
compile:
	go get github.com/theaayushpatel/BlackJack/... # partially downloads again, but we need the library packages in go path to build
	# <--- second compilation, maybe -d flag on go get above is better
	env GOBIN=$(CURDIR)/build go install ./cmd/... # compile all main packages to the build folder

	# compile the web app
	# ToDo: (check if dependencies have been fetched by 'go get', even with the build js tags)
	$(HOME)/go/bin/gopherjs get github.com/theaayushpatel/BlackJack/web_client/...
	$(HOME)/go/bin/gopherjs build -m -o build/webapp.js web_client/*.go

installkali:
	#apt-get -y install git screen hostapd autossh bluez bluez-tools bridge-utils policykit-1 genisoimage iodine haveged
	#apt-get -y install tcpdump
	#apt-get -y install python-pip python-dev

	# before installing dnsmasq, the nameserver from /etc/resolv.conf should be saved
	# to restore after install (gets overwritten by dnsmasq package)
	#cp /etc/resolv.conf /tmp/backup_resolv.conf
	#apt-get -y install dnsmasq
	#/bin/bash -c 'cat /tmp/backup_resolv.conf > /etc/resolv.conf'

	# python dependencies for HIDbackdoor
	sudo pip install pydispatcher

	cp build/BlackJack_service /usr/local/bin/
	cp build/BlackJack_cli /usr/local/bin/
	cp dist/BlackJack.service /etc/systemd/system/BlackJack.service
	# copy over keymaps, scripts and www data
	mkdir -p /usr/local/BlackJack
	cp -R dist/keymaps /usr/local/BlackJack/
	cp -R dist/scripts /usr/local/BlackJack/
	cp -R dist/HIDScripts /usr/local/BlackJack/
	cp -R dist/www /usr/local/BlackJack/
	cp -R dist/db /usr/local/BlackJack/
	cp -R dist/helper /usr/local/BlackJack/
	cp -R dist/ums /usr/local/BlackJack/
	cp -R dist/legacy /usr/local/BlackJack/
	cp build/webapp.js /usr/local/BlackJack/www
	cp build/webapp.js.map /usr/local/BlackJack/www

	# careful testing
	#sudo update-rc.d dhcpcd disable
	#sudo update-rc.d dnsmasq disable
	systemctl disable networking.service # disable network service, relevant parts are wrapped by BlackJack (boottime below 20 seconds)

	# enable service
	systemctl enable haveged
	systemctl enable avahi-daemon
	systemctl enable BlackJack.service

install:
	cp build/BlackJack_service /usr/local/bin/
	cp build/BlackJack_cli /usr/local/bin/
	cp dist/BlackJack.service /etc/systemd/system/BlackJack.service
	# copy over keymaps, scripts and www data
	mkdir -p /usr/local/BlackJack
	cp -R dist/keymaps /usr/local/BlackJack/
	cp -R dist/scripts /usr/local/BlackJack/
	cp -R dist/HIDScripts /usr/local/BlackJack/
	cp -R dist/www /usr/local/BlackJack/
	cp -R dist/db /usr/local/BlackJack/
	cp dist/bin/* /usr/local/bin/
	cp build/webapp.js /usr/local/BlackJack/www
	cp build/webapp.js.map /usr/local/BlackJack/www

	# careful testing
	#sudo update-rc.d dhcpcd disable
	#sudo update-rc.d dnsmasq disable
	# systemctl disable networking.service # disable network service, relevant parts are wrapped by BlackJack (boottime below 20 seconds)

	# reinit service daemon
	# systemctl daemon-reload
	# enable service
	# systemctl enable haveged
	# systemctl enable BlackJack.service
	# start service
	# service BlackJack start

remove:
	# stop service
	service BlackJack stop
	# disable service
	systemctl disable BlackJack.service
	rm -f /usr/local/bin/BlackJack_service
	rm -f /usr/local/bin/BlackJack_cli
	rm -f /etc/systemd/system/BlackJack.service
	rm -R /usr/local/BlackJack/    # this folder should be kept, if only an update should be applied
	# reinit service daemon
	systemctl daemon-reload

	#sudo update-rc.d dhcpcd enable

