check_install:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	swagger	generate spec -o ./swagger.json --scan-models