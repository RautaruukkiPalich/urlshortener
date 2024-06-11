run:
	go run cmd/urlsh/main.go --config=./config/config.yaml

rundocker:
	docker-compose up -d

stopdocker:
	docker-compose down

test: testunit testintegration

testunit: 
	go test ./... -v

testintegration:
	go test ./test/... -v

lint:
	golangci-lint run ./...

swagger:
	swag init -g cmd/urlsh/main.go
	swag fmt

bomb:
	bombardier -c 1000 -d 30s -l http://localhost:8080/05JrT5_c

	bombardier -c 1000 -d 30s http://localhost:8080/t/05JrT5_c
