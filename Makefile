test:
	@go test -v

cover:
	@go test -coverprofile=c.out && go tool cover -html=c.out