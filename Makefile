export GOPATH=${PWD}/.go

deps:
	 go get -d

build_module: deps
	go build -buildmode=c-shared -o pam_authy.so

build_helper: deps
	go build -o pam_authy_helper

install:
	cp pam_authy.so /lib/x86_64-linux-gnu/security/

test: install
	pamtester sshd root authenticate
