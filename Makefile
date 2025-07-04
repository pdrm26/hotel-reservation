build:
	@go build -o bin/api

run: build
	@./bin/api

test:
	@go test -v ./...

seed:
	@go run scripts/seed.go

docker:
	echo "building docker file"
	@docker build -t api .
	@docker run -p 3000:3000 api
