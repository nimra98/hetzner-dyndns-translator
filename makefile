all: build build-docker

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./hetzner-dyndns-translator codeberg.org/anbraten/hetzner-dyndns-translator

build-docker:
	docker build --rm -t anbraten/hetzner-dyndns-translator:latest .

release: build build-docker
	docker push anbraten/hetzner-dyndns-translator:latest
