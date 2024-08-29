# Hetzner Dyndns Translator

Hetzner Dyndns Translator ermöglicht den Zugriff auf die DNS API von Hetzner in einem Format, das mit dem bei Fritz!Boxen eingebauten DynDNS-Dienst kompatibel ist.
Die Software ist als Docker-Image verfügbar und kann einfach über Docker Hub bezogen werden.

Hetzner Dyndns Translator provides access to the Hetzner DNS API in a format compatible with the built-in DynDNS service of Fritz!Box devices.
The software is available as a Docker image and can be easily obtained from Docker Hub.

## Docker Hub

The docker images are available on [Docker Hub](https://hub.docker.com/r/nimra98/hetzner-dyndns-translator).

## Usage with docker (compose)

### Usage with docker CLI

```bash
docker run --rm -p 3000:3000 -e SERVICE_AUTH_TOKEN=mysupersecrettoken -e SHOW_HETZNER_TOKEN=true nimra98/hetzner-dyndns-translator:latest
```

### Standalone http servicce
```yaml
services:
  translator:
    image: nimra98/hetzner-dyndns-translator:latest
    environment:
      - PORT=3000 # optional, default is 3000, the port the server listens on
      - SERVICE_AUTH_TOKEN=mysupersecrettoken # optional, default is none, the token that is required to access the service
      - SHOW_HETZNER_API_TOKEN=true # optional, default is false, if set to true the hetzner api token is shown in the logs
    ports:
      - 3000:3000 # Depends on the PORT environment variable above
```

### Behind traefik reverse proxy
```yaml
services:
  dyndns-translator:
    image: nimra98/hetzner-dyndns-translator:latest
    container_name: dyndns-translator
    restart: always
    labels:
      - "traefik.enable=true"
      - "traefik.docker.network=web-proxy"
      - "traefik.http.routers.dyndns-translator-ondomain.middlewares=sec@file, gzip@file"
      - "traefik.http.routers.dyndns-translator-ondomain.rule=Host(`dyndns-translator.ondomain.tld`)"
      - "traefik.http.routers.dyndns-translator-ondomain.tls.options=intermediate@file"
      - "traefik.http.routers.dyndns-translator-ondomain.tls.certresolver=httpchallenge"
      - "traefik.http.services.dyndns-translator.loadbalancer.server.port=3000" # Depends on the PORT environment variable above
    environment:
        - PORT=3000 # optional, default is 3000, the port the server listens on
        - SERVICE_AUTH_TOKEN=mysuperscrettoken # optional, default is none, the token that is required to access the service
        - SHOW_HETZNER_API_TOKEN=true # optional, default is false, if set to true the hetzner api token is shown in the logs
    networks:
      - default

networks:
  default:
    external:
      name: web-proxy # The name of the traefik network
```

## Update records

```bash
# update A record
curl dyndns-translator.ondomain.tld[:Port]/[SERVICE_AUTH_TOKEN]/dyndns/subdomainwithoutzonepart/example.tld/Hetzner_API_Token/$(curl -s http://v4.ipv6-test.com/api/myip.php)

# update AAAA record
curl dyndns-translator.ondomain.tld[:Port]/[SERVICE_AUTH_TOKEN]/dyndns/subdomainwithoutzonepart/example.tld/Hetzner_API_Token/$(curl -s http://v6.ipv6-test.com/api/myip.php)
```

## Build and push translator server

```bash
# Build the docker image, tag it with the version and latest, store in the local registry
sudo make build VERSION=1.0.0 LATEST=true

# Build the docker image, tag it with the version and latest, push it to the docker hub
sudo make release VERSION=1.0.0 LATEST=true
```