.PHONY: build clean deploy

clean:
	rm -rf ./output ./.serverless

build: clean
	# Create output directories
	mkdir -p ./output/handlers
	mkdir -p ./output/tmp

	# Build executables for Lambda
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./output/tmp/create handlers/create.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./output/tmp/delete handlers/delete.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./output/tmp/get handlers/get.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./output/tmp/list handlers/list.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./output/tmp/update handlers/update.go

	# Create deployment packages (zip files)
	cp ./output/tmp/create ./output/tmp/bootstrap && (cd ./output/tmp && zip -m ../handlers/create.zip bootstrap)
	cp ./output/tmp/delete ./output/tmp/bootstrap && (cd ./output/tmp && zip -m ../handlers/delete.zip bootstrap)
	cp ./output/tmp/get ./output/tmp/bootstrap && (cd ./output/tmp && zip -m ../handlers/get.zip bootstrap)
	cp ./output/tmp/list ./output/tmp/bootstrap && (cd ./output/tmp && zip -m ../handlers/list.zip bootstrap)
	cp ./output/tmp/update ./output/tmp/bootstrap && (cd ./output/tmp && zip -m ../handlers/update.zip bootstrap)

deploy: build
	sls deploy --verbose
