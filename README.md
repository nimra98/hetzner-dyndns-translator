# Hetzner Dyndns Translator

## Run translator server

```bash
docker run --rm -ti --port 3000:3000 nimra98/hetzner-dyndns-translator:latest
```

## Update records

```bash
# update A record
curl my-server.tld:3000/dyndns/subdomainwithoutzonepart/example.tld/mysupersecrettoken/$(curl -s http://v4.ipv6-test.com/api/myip.php)

# update AAAA record
curl my-server.tld:3000/dyndns/subdomainwithoutzonepart/example.tld/mysupersecrettoken/$(curl -s http://v6.ipv6-test.com/api/myip.php)
```
