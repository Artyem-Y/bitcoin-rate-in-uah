# Makefile

include $(PWD)/.env
export

run:
	docker build . -t bitcoin-rate-in-uah
	docker run -v "${PWD}/:/app" -p 8080:8080 --rm --env-file=".env" bitcoin-rate-in-uah

start:
	go run main.go