test:
	go test -v -coverprofile cover.out

coverage:
	go test
	go tool cover -html=cover.out
