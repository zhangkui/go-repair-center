# Benzhi Docker Build

This repository builds the Go Repair Center API from `cmd/api`.

```sh
DOCKER_PLATFORM=linux/amd64 IMAGE_NAME=go-repair-center ./build_benzhi_docker.sh
```

The image contains the API binary, migrations, and the example environment file.
