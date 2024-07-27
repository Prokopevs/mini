build_image:
	docker build -t server/message-api:latest -f ./server/Dockerfile ./server