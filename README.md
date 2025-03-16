# Mattress

*"Hide your cash under the mattress!"*

Nostr relay for Cashu tokens of whitelisted npubs.

Based on and inspired by the [haven](https://github.com/bitvora/haven) relay.

## Build

Build and run it locally with go:
```bash
# rm -rf db
# go mod download
# go mod tidy
go build
./mattress
```

Or use a container (podman):
```bash
podman build -t mattress:test .
podman run -it --rm --replace --name mattress -p 127.0.0.1:2121:2121 -v ./.env:/app/.env:Z -v ./user_npubs.json:/app/user_npubs.json:Z localhost/mattress:test
```
