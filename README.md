# Helm EchoServer

A simple Go echoserver built with [Echo](https://github.com/labstack/echo) framework.

## Endpoints
- `GET /` prints information about Kubernetes, Helm, K8S Pod and HTTP request/response. 
- `GET /env` prints env variables in the environemtn where the server is running.
- `GET /content` reads the content of /tmp/content in the environemtn where the server is running.
- `POST /content` appends to the content of /tmp/content in the environemtn where the server is running.

## Run
> default port is 8080, but can be overriden 

- From source: `go run server.go -p <server port>`
- Docker: `docker run -e PORT=<server port> -p <host port>:<container port> praqma/helm-echoserver:1.0`

> docker e.g. docker run -e PORT=80 -p 8080:80 praqma/helm-echoserver:1.0 
> Then access it on http://localhost:8080
