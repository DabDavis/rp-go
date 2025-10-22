# rp-go

Game engine initialized with a modular ECS + Pub/Sub architecture.

## Testing

Simulation and rendering are now fully separated, allowing the engine to run
its logic systems without a GPU or windowing system. The platform layer ships
with a pure Go headless implementation that satisfies the same API surface as
the Ebiten-backed runtime, so tests can execute deterministically inside CI
containers.

The helper script automatically applies the `headless` build tag:

```bash
./run-tests.sh            # equivalent to: go test -tags headless ./...
./run-tests.sh ./engine/  # passes through additional go test arguments
```

You can also invoke `go test -tags headless ./...` directly if you prefer to
manage the invocation yourself.
