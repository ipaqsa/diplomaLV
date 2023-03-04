node-build:
	go build -o build/node ./cmd/main.go

node-run: node-build
	build/node -c config/node.yml