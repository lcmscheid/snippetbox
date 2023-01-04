# To generate self-signed certificates using go standard library tools
1. install go
2. check the go environment to get the path to the go src
```bash
$ go env | grep GOROOT 
```
3. create a folder in the root of the project to hold self-signed certificates
```bash
$ mkdir -p tls & cd tls
```
4. Generate self-signed certificate
```bash
$ go run <GOROOT>/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```