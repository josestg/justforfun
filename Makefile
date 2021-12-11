
help: ## prints make help.
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

hooks: ## register git-hook scripts.
	@chmod +x vars/scripts/pre-commit
	@git config core.hooksPath vars/scripts

test: ## executes test files. (race_detector: enabled, cache: disable)
	go test -race -v -count=1 ./...

httpd: ## compiles httpd into standalone binary.
	go build \
		-ldflags "-X main.buildName=httpd-justforfun \
				  -X main.buildRef=`git rev-parse --short HEAD` \
				  -X main.buildDate=`date -u +"%Y-%m-%dT%H:%M:%SZ"`" \
		-o httpd.so cmd/httpd/main.go

httpd-image: ## builds a docker image of httpd application.
	docker build \
		-f cmd/httpd/Dockerfile \
		-t josestg/justforfun-httpd \
		--build-arg IMAGE_NAME=josestg/justforfun-httpd \
		--build-arg BUILD_REF=`git rev-parse --short HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg VENDOR=josestg \
		--build-arg SOURCE=github.com/josestg/justforfun/cmd/httpd/Dockerfile \
		--build-arg AUTHOR="Jose Sitanggang <josealfredositanggang@gmail.com>" \
		.