test:
	go test -v -coverprofile cover.out

coverage: test 
	go tool cover -html=cover.out
