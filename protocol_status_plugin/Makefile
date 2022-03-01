build:
	@go build

run:
	@go run main.go

coverage:
	go test -coverprofile cover.out

reflect:
	@grpcurl -plaintext localhost:9091 list pilot.plugin.ManagerPlugin

mock:
	mockery -r --name Executor

clean:
	rm -rf cover.out
	rm -rf vatz-plugin-matic
