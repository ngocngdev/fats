#!/bin/bash
source ./util.sh
set -x


if [[ "$OS_NAME" == "windows" ]]; then
	choco install kubernetes-cli
	mkdir -p ~/.kube
else
	curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.9.2/bin/linux/amd64/kubectl
	chmod +x kubectl
	sudo mv kubectl /usr/local/bin/
fi


if [[ "$OS_NAME" == "windows" ]]; then
	go get -d github.com/boz/kail
	cd $GOPATH/src/github.com/boz/kail
	# temporarily mimics what the kail's Makefile does
	go get -u github.com/kardianos/govendor
	$GOPATH/govendor sync -v
	$GOPATH/govendor build -i +program
	mv kail.exe /usr/local/bin/kail
else
	mkdir -p kail
	curl -L https://github.com/boz/kail/releases/download/v0.6.0/kail_0.6.0_linux_amd64.tar.gz \
	  | tar xz -C kail
	chmod +x kail/kail
	sudo mv kail/kail /usr/local/bin/
fi
