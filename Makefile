run: rundocker
	go run cmd/urlsh/main.go --config=./config/config.yaml

rundocker:
	docker-compose up -d

test: testunit testintegration

testunit: 
	go test ./... -cover

testintegration:
	go test ./test/...

lint:
	golangci-lint run ./...

swagger:
	swag init -g cmd/urlsh/main.go
	swag fmt