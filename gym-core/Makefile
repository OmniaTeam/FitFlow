include .env
export

diagram:
	@atlas schema inspect --env server -w

up:
	air -c .air.toml

swagger-generate:
	@swag init -g ./cmd/omnia/web.go

build:
	go build -o tmp/core ./cmd/omnia/web.go

pull:
	git pull

restart: pull build
	systemctl restart omnia_core.service

hash:
	@atlas migrate hash

migrate: hash
	@atlas migrate apply --env server

logs:
	@journalctl -u omnia_core.service -f -n 50

deploy:
	ssh root@147.45.154.109 'make -C /home/omnia/gym-core restart'