# References

https://github.com/dcu/go-authy - Authy API in Golang
https://github.com/uber/pam-ussh - PAM module in Golang
https://github.com/authy/authy-ssh - Motivation (do it better)
https://github.com/golang/dep/ - why TOML instead of YAML or JSON

# Testing

```
./build1
./build2
docker run --rm -ti -v $(pwd):/app -w /app -p 2222:22 testsshd2 bash
make test
```

# TODO

* CI builds and tests
* Specify different conf file path
* cmd utility to register users