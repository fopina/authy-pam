all:
	go build -buildmode=c-shared -o pam_authy.so

make install: all
	cp pam_authy.so /lib/x86_64-linux-gnu/security/
