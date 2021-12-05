
help: ## prints make help.
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

httpd: ## compiles httpd into standalone binary.
	go build \
		-ldflags "-X main.buildName=httpd-justforfun \
				  -X main.buildRef=`git rev-parse --short HEAD` \
				  -X main.buildDate=`date -u +"%Y-%m-%dT%H:%M:%SZ"`" \
		-o httpd.so cmd/httpd/main.go