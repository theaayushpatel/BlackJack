FROM kalilinux/kali-rolling

WORKDIR /root
RUN apt-get update && apt-get -y install git wget nano
# install Go 1.12 instead of Kali bundled Go 1.13 (GopherJS needs 1.12)
RUN wget https://dl.google.com/go/go1.12.16.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.12.16.linux-amd64.tar.gz
ENV PATH "$PATH:/usr/local/go/bin:/root/go/bin"
RUN go get -u github.com/gopherjs/gopherjs
# clone BlackJack master (has to be changed in order to use a different branch/tag)
RUN git clone https://github.com/mame82/BlackJack_aloa

# BlackJack webclient dependencies
RUN go get -u github.com/johanbrandhorst/protobuf/...
# manual population of go source tree, with dependencies of BlackJack webclient
# using git clone.
# This is really messy, but unfortunately GopherJS has no module support (dependency tracking)
# At least this leaves room to modify "git clone" commands to grab proper branches,
# in case this is required
RUN mkdir -p /usr/local/go/src/github.com/mame82/
RUN git clone https://github.com/mame82/hvue /usr/local/go/src/github.com/mame82/hvue
RUN git clone https://github.com/mame82/mvuex /usr/local/go/src/github.com/mame82/mvuex
# copy already cloned repo of BlackJack instead of cloning (assures same branch in Go
# source tree)
RUN cp -R BlackJack_aloa/ /usr/local/go/src/github.com/mame82/


# run a test build, otherwise the image could be used interactively 
# with build.sh as used below
WORKDIR /root/BlackJack_aloa/build_support
RUN ./build.sh


