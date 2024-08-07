genpb:
	go run ./tool/goit-protoc/*.go ./app/hellogrpc/proto/

airweb:
	air --build.cmd "go build -o ./tmp/main ./cmd/helloweb-server"

at:
	go test ./acceptancetest/helloweb -run TestWeb -update

.PHONY: genpb