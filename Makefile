.PHONY: build clean deploy

clean:
	rm -rf ./output ./.serverless Gopkg.lock

build: clean
	# Build executables and place them directly at the root of each function directory
	cd get && env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o ../output/get/bootstrap get.go && cd ..
	cd get && env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o ../output/getquery/bootstrap getQuery.go && cd ..
	cd post && env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o ../output/post/bootstrap post.go && cd ..

	# Create deployment packages (zip files)
	cd ./output/get && zip -r get.zip bootstrap
	cd ./output/getquery && zip -r getquery.zip bootstrap
	cd ./output/post && zip -r post.zip bootstrap

deploy: build
	sls deploy --verbose
