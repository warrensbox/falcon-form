build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/send-msg main.go

publish:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/message-falcon main.go
	zip message-falcon.zip bin/*
	aws s3 cp message-falcon.zip s3://message-falcon