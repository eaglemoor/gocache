test:
	@go test -v -coverprofile=c.out -covermode=count

cover:
	@go test -coverprofile=c.out -covermode=count && go tool cover -html=c.out