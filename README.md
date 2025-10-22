# rp-go

Game engine initialized with a modular ECS + Pub/Sub architecture.

## Testing

The test suite exercises systems that require an active display context. When
running in headless environments (such as CI) install `xvfb` so the helper
script can provision a virtual framebuffer automatically. The script prefers
`xvfb-run` when available and otherwise boots a background `Xvfb` server.

```bash
sudo apt-get update && sudo apt-get install -y xvfb
./run-tests.sh
```

`run-tests.sh` forwards any additional flags to `go test` if you need to target a
specific package, e.g.

```bash
./run-tests.sh ./engine/...
```
