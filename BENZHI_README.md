# Benzhi Docker Build

This repository builds the Go Repair Center API from cmd/api.

Build with:

`
DOCKER_PLATFORM=linux/amd64 IMAGE_NAME=go-repair-center ./build_benzhi_docker.sh
`

The image is built from enzhi.Dockerfile and includes the API binary, migrations, and the example environment file.
