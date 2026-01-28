build:
	@go build -o ./dist/app
run: build
	@./dist/app	
test: 
	@go clean -testcache
	@go test -v ./ 
test-race: 
	@go clean -testcache
	@go test -race -v ./ 
