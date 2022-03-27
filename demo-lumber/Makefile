up: data/nginx
	docker compose up -d filebeat nginx-a nginx-b
	docker compose logs -f filebeat

pull:
	docker compose pull

down:
	docker compose down --remove-orphans

client:
	docker compose run client

status:
	docker compose ps

data/nginx:
	mkdir -p data/nginx
