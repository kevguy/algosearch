# AlgoSearch

## Folder Structure

```
- backend
- vendor
- zarf
- dockerfile.algosearch-backend.dockerignore
- dockerfile.metrics.dockerignore
- go.mod
- makefile
- README.md
```

- backend: contains all the Go application code that's required to run the backend api
- vendor: contains all the Go dependencies the project needs to run the backend
- zarf: contains all the Docker and Docker compose related files
- dockerfile.algosearch-backend.dockerignoe: the `.dockerignore` file needed for building the `algosearch-backend` image
- dockerfile.metrics.dockerignoe: the `.dockerignore` file needed for building the `metrics` image
- go.mod: contains the list of Go dependencies the backend needs
- makefile: contains a general list of useful commands to run the project
