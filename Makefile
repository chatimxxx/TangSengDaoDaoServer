build:
	docker build -t xochat_im_server .
push:
	docker login
	docker tag xochat_im_server zyuan909/xochat_im_server:latest
	docker push zyuan909/xochat_im_server:latest
deploy:
	docker login
	docker build -t xochat_im_server .
	docker tag xochat_im_server zyuan909/xochat_im_server:latest
	docker push zyuan909/xochat_im_server:latest
deploy-arm:
	docker login
	docker build -t xochat_im_serverarm64 .
	docker tag xochat_im_serverarm64 zyuan909/xochat_im_server:latest-arm64
	docker push zyuan909/xochat_im_server:latest-arm64
deploy-v1.5:
	docker build -t xochat_im_server .
	docker tag xochat_im_server zyuan909/xochat_im_server:v1.5
	docker push zyuan909/xochat_im_server:v1.5
run-dev:
	docker-compose build;docker-compose up -d
stop-dev:
	docker-compose stop
env-test:
	docker-compose -f ./testenv/docker-compose.yaml up -d 