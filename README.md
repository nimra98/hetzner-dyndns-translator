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
docker run --rm -p 3000:3000 -e SERVICE_AUTH_TOKEN=mysupersecrettoken -e SHOW_HETZNER_TOKEN=true -e HETZNER_API_VERSION=cloud nimra98/hetzner-dyndns-translator:latest
```

### Configuration Parameters (Environment Variables)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | The port the server listens on. |
| `SERVICE_AUTH_TOKEN` | `none` | Optional token required to access the translator service. |
| `SHOW_HETZNER_API_TOKEN` | `false` | If set to `true`, the Hetzner API token is shown in the logs (use for debugging only!). |
| `HETZNER_API_VERSION` | `legacy` | Switch between `legacy` (dns.hetzner.com) and `cloud` (api.hetzner.cloud) API. |

### API Migration (Hetzner Console)

Hetzner is migrating its DNS service to the main Hetzner Console (Cloud API). This tool supports both:

1. **Legacy API (`HETZNER_API_VERSION=legacy`):** Uses tokens from `dns.hetzner.com`.
2. **Cloud API (`HETZNER_API_VERSION=cloud`):** Uses project-based API tokens from the new Hetzner Console. **Required for migrated zones!**

### Docker Compose Example

```yaml
services:
  dyndns-translator:
    image: nimra98/hetzner-dyndns-translator:latest
    container_name: dyndns-translator
    restart: always
    environment:
      - PORT=3000
      - SERVICE_AUTH_TOKEN=mysupersecrettoken
      - HETZNER_API_VERSION=cloud # Set to 'cloud' for the new Hetzner DNS Console
      - SHOW_HETZNER_API_TOKEN=false
    ports:
      - 3000:3000
```

## Update records

The update URL format remains identical for both API versions. However, the `Hetzner_API_Token` must match the configured `HETZNER_API_VERSION`.

```bash
# update A record
curl dyndns-translator.ondomain.tld[:Port]/dyndns/[SERVICE_AUTH_TOKEN]/subdomainwithoutzonepart/example.tld/Hetzner_API_Token/$(curl -s http://v4.ipv6-test.com/api/myip.php)

# update AAAA record
curl dyndns-translator.ondomain.tld[:Port]/dyndns/[SERVICE_AUTH_TOKEN]/subdomainwithoutzonepart/example.tld/Hetzner_API_Token/$(curl -s http://v6.ipv6-test.com/api/myip.php)
```

For the Fritz!Box configuration, the following values are required:

| Setting               | Value                                      |
|-----------------------|--------------------------------------------|
| Update URL            | `dyndns-translator.ondomain.tld[:Port]/dyndns/[SERVICE_AUTH_TOKEN]/subdomainwithoutzonepart/example.tld/Hetzner_API_Token/<ip6addr>`    |
| Domain name           | Does not matter                            |
| Username              | Does not matter                            |
| Password              | Does not matter                            |

## Build and push translator server

The project now requires **Go 1.25** for building.

```bash
# Build the docker image, tag it with the version and latest, store in the local registry
make build VERSION=1.0.0 LATEST=true

# Build the docker image, tag it with the version and latest, push it to the docker hub
make release VERSION=1.0.0 LATEST=true
```

## Known Limitations (Hetzner Cloud API)

When using `HETZNER_API_VERSION=cloud`, please be aware of the following technical details:

1. **Asynchronous Updates:** The new Hetzner Cloud API processes DNS changes as asynchronous "Actions". This tool considers an update successful as soon as the API accepts the request (`201 Created`). It does **not** wait for the action to finish (`finished: null`). In rare cases, the DNS update might take a few seconds to become active.
2. **Pagination:** Currently, the tool only fetches the first page of DNS zones (default limit is usually 25). If you manage a large number of zones, ensure the target zone is among the first 25, or contribute a pagination fix.
3. **Record Matching:** The tool identifies records by their name and type (`A` or `AAAA`). Ensure that the record exists in the Hetzner Console before the first update.
