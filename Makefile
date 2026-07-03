PROTO_DIR=internal/gateway/grpc/inventory/pb

proto:
	cd $(PROTO_DIR) && PATH="$$PATH:$$(go env GOPATH)/bin" buf dep update && PATH="$$PATH:$$(go env GOPATH)/bin" buf generate

docker_build:
	docker build --no-cache -t reservation:latest .

docker_up:
	docker compose up -d

docker_down:
	docker compose down --volumes
	docker container prune -f
	docker image prune -f
	docker volume prune -f

docker_restart: docker_down docker_build docker_up
docker_start: docker_down docker_up



k8s-up:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secret.yaml
	kubectl apply -f k8s/postgres.yaml
	kubectl apply -f k8s/redis.yaml
	kubectl apply -f k8s/api.yaml
	kubectl apply -f k8s/grpc.yaml

k8s-down:
	kubectl delete -f k8s/

k8s-status:
	kubectl get all -n inventory