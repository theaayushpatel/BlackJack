#!/bin/bash
echo compiling ...
env GOOS=linux GOARCH=arm GOARM=6 go build -gcflags "all=-N -l" -o build/BlackJack_service cmd/BlackJack_service/BlackJack_service.go
env GOOS=linux GOARCH=arm GOARM=6 go build -o build/BlackJack_cli cmd/BlackJack_cli/BlackJack_cli.go
env GOOS=linux GOARCH=arm GOARM=6 go build -o /tmp/ntest ntest.go

echo uploading ...
scp /tmp/ntest 172.16.0.1:~/BlackJack/build
scp build/BlackJack_service 172.16.0.1:~/BlackJack/build
scp build/BlackJack_cli 172.16.0.1:~/BlackJack/build

