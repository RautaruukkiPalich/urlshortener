run:
	go run cmd/urlsh/main.go --config=./config/config.yaml

rundocker:
	docker-compose up -d

stopdocker:
	docker-compose down


test: testunit testintegration

testunit: 
	go test ./... 

testintegration:
	go test ./test/... 

lint:
	golangci-lint run ./...

swagger:
	swag init -g cmd/urlsh/main.go
	swag fmt

bomb:
	bombardier -c 200 -d 30s -l http://localhost:8080/dg4r-TMc

	bombardier -c 1000 -d 30s http://localhost:8080/1
