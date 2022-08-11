# public-api-design

This repository contains Goa designs for Flexera's public APIs.

Related Resources


## Project Structure

```shell
GOPATH
├─src ...
├─├─github.com                    # Test files (alternatively `spec` or `tests`)
│   ├── mlb          # Load and stress tests
│   │   ├── mlb-ballpark-segregation-service         # End-to-end, integration tests (alternatively `e2e`)
│   │   │   └── front-service/design
|   |   |   └── DockerFile        # Unit tests
└── ...

```

## Goa Generate 

Run the below command to generate design files

```shell
goa gen github.com/mlb/mlb-ballpark-segregation-service/front_service/design
```
