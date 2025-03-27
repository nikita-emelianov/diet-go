.PHONY: build clean deploy

clean:
	rm -rf ./output ./.serverless

build: clean
	# Create output directories
	mkdir -p ./output/handlers
	mkdir -p ./output/tmp

	# Build executables and place them directly at the root of each function directory
	env GOARCH=arm64 GOOS=linux
	go build -ldflags="-s -w" -o ./output/tmp/get handlers/get.go
	go build -ldflags="-s -w" -o ./output/tmp/get_query handlers/get_query.go
	go build -ldflags="-s -w" -o ./output/tmp/post handlers/post.go

	# Create deployment packages (zip files)
	cp ./output/tmp/get ./output/tmp/bootstrap && (cd ./output/tmp && zip -m ../handlers/get.zip bootstrap)
	cp ./output/tmp/get_query ./output/tmp/bootstrap && (cd ./output/tmp && zip -m ../handlers/get_query.zip bootstrap)
	cp ./output/tmp/post ./output/tmp/bootstrap && (cd ./output/tmp && zip -m ../handlers/post.zip bootstrap)

deploy: build
	sls deploy --verbose
