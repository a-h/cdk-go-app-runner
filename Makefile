build-linux:
	cd app && GOOS=linux GOARCH=amd64 go build -o ../docker-images/app

deploy: build-linux
	npx cdk deploy
