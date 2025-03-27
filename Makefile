.PHONY: build clean deploy

clean:
	rm -rf ./output ./.serverless

build: clean
	# Build executables and place them directly at the root of each function directory
	env GOARCH=arm64 GOOS=linux
	go build -ldflags="-s -w" -o output/bootstrap/get/bootstrap handlers/get.go
	go build -ldflags="-s -w" -o output/bootstrap/get_query/bootstrap handlers/get_query.go
	go build -ldflags="-s -w" -o output/bootstrap/post/bootstrap handlers/post.go

	# Create deployment packages (zip files)
	mkdir -p ./output/handlers
	zip -r ./output/handlers/get.zip ./output/bootstrap/get/bootstrap
	zip -r ./output/handlers/get_query.zip ./output/bootstrap/get_query/bootstrap
	zip -r ./output/handlers/post.zip ./output/bootstrap/post/bootstrap

deploy: build
	sls deploy --verbose
