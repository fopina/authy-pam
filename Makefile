all: pam_authy.so
	GOPATH=${PWD}/.go go build -buildmode=c-shared -o pam_authy.so

install: all
	cp pam_authy.so /lib/x86_64-linux-gnu/security/

test: install
	pamtester sshd root authenticate
