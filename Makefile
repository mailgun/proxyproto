help: # Display help
	@awk -F ':|##' \
			'/^[^\t].+?:.*?##/ {\
					printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
			}' $(MAKEFILE_LIST)

run:
	go run ./cmd/proxy-client/main.go

.PHONY: run
