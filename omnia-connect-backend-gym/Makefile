.PHONY: logs
up:
	air -c .air.toml

swagger-generate:
	swag init -g ./cmd/sso/main.go

build:
	go build -o tmp/sso ./cmd/sso/main.go

pull:
	git pull

restart: pull build
	systemctl restart sso-omnia.service

logs:
	journalctl -u sso-omnia.service -f -n 50

deploy:
	ssh root@147.45.154.109 'make -C /home/sso/omnia-connect-backend restart'