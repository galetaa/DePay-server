# DevOps

The repository supports two local modes.

## Dev Mode

Dev mode uses `go run` with source mounts:

```bash
make dev-ready
make full-up
```

## Image Mode

Image mode builds production-like images with service Dockerfiles and starts them without source mounts:

```bash
make docker-build
make image-up
make image-down
```

`docker-compose.images.yml` expects the shared infrastructure services from `docker-compose.yml` (`postgres`, `redis`, `rabbitmq`) to run in the same compose project.

## CI

`.github/workflows/ci-cd.yml` includes:

- Go tests per module
- web tests and build
- SQL migrations, seed, and SQL tests against PostgreSQL
- Docker image build
- compose health smoke
- Kubernetes manifest client validation

## Kubernetes

Kubernetes manifests live in `k8s/`. `make k8s-client-validate` performs offline YAML parsing without requiring a cluster; `make k8s-validate` keeps the stricter server-side validation path for a configured cluster.
