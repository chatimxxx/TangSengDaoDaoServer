build:
	docker build -t tsdd_server .
push:
	docker login
	docker tag tsdd_server zyuan909/tsdd_server:latest
	docker push zyuan909/tsdd_server:latest
deploy:
	docker login
	docker build -t tsdd_server .
	docker tag tsdd_server zyuan909/tsdd_server:v1.0.2
	docker push zyuan909/tsdd_server:v1.0.2
deploy-v1.5:
	docker build -t tsdd_server .
	docker tag tsdd_server zyuan909/tsdd_server:v1.5
	docker push zyuan909/tsdd_server:v1.5
run-dev:
	docker-compose build;docker-compose up -d
stop-dev:
	docker-compose stop
env-test:
	docker-compose -f ./testenv/docker-compose.yaml up -d 