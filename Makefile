build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/message-falcon main.go

publishcircle:
	go get -v -t -d ./...
	env GOOS=linux go build -ldflags="-s -w" -o bin/message-falcon main.go
	zip message-falcon-${CIRCLE_BUILD_NUM}.zip bin/*
	aws s3 cp message-falcon-${CIRCLE_BUILD_NUM}.zip s3://message-falcon
	aws lambda update-function-code --function-name message-falcon --s3-bucket message-falcon --s3-key message-falcon-${CIRCLE_BUILD_NUM}.zip

publish:
	go get -v -t -d ./...
	env GOOS=linux go build -ldflags="-s -w" -o bin/message-falcon main.go
	zip message-falcon.zip bin/*
	aws s3 cp message-falcon.zip s3://message-falcon
	aws lambda update-function-code \
		--function-name  message-falcon \
		--zip-file fileb://./message-falcon.zip

	

upload:
	aws s3 cp --recursive --acl public-read docs/_site s3://falcon-form.warrensbox.com/ 
