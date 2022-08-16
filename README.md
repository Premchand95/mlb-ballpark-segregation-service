# MLB-Ballpark-services

The MLB BallPark Segregation Service API provides REST-ful API to retrieve custom schedules of mlb games

## Prerequisites

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/) [![Docker](https://badgen.net/badge/icon/docker?icon=docker&label)](https://https://docker.com/) [![Linux](https://svgshare.com/i/Zhy.svg)](https://svgshare.com/i/Zhy.svg) [![git](https://badgen.net/badge/icon/git?icon=git&label)](https://git-scm.com)
* [go1.18.2](https://go.dev/dl/)
* [docker 20.10.17](https://docs.docker.com/get-docker/)
* [goa v3.7.5](https://github.com/goadesign/goa)
* [mockgen v1.6.0](https://github.com/golang/mock)

### set up
 clone the repo to $GOPATH/src/github.com/mlb folder & make sure we have all Prerequisites in our machine
 ```shell
 git clone https://github.com/Premchand95/mlb-ballpark-segregation-service.git
 ```

## Project Structure

```shell
GOPATH
├─src ...
├─├─github.com                    
│   ├── mlb          
│   │   ├── mlb-ballpark-segregation-service
│   │   │   └── front-service/design
|   |   |   └── front-service/gen
|   |   |   └── Makefile
|   |   |   └── DockerFile       
└── ...

```
## Run API with docker & make
    
  ### Build Docker Image 

  run the below command in mlb-ballpark-segregation-service folder to generate docker image `mlb-segregation-service`

  ```shell
  make docker-image
  ```
  see docker images in machine by running below command

   ```
   docker images
   ```
   output:
   ```shell
      REPOSITORY                TAG             IMAGE ID       CREATED          SIZE
      mlb-segregation-service   latest          e07c9e81e3cb   29 minutes ago   21.4MB
      <none>                    <none>          f9d21a4fbbc6   29 minutes ago   458MB
      alpine                    latest          9c6f07244728   6 days ago       5.54MB
      golang                    1.18.2-alpine   5432f005be2d   2 months ago     328MB
   ```
   Run the below command to run the docker container with image `mlb-segregation-service` exposing port 8080
   ```shell
   make run
  ```
  output:
  ```shell
      docker run -p 8080:8080 mlb-segregation-service:latest
      [MLB BallPark Segregation Service API] 03:44:27 HTTP "Index" mounted on GET /api/v1/teams/{id}/schedule
      [MLB BallPark Segregation Service API] 03:44:27 HTTP server listening on "0.0.0.0:8080"
  ```
  ### Run API without docker
  
  Just run the below command and start testing
  
  ```shell
  go run front_service/cmd/front_service/*
  ```
  
  ### Make call to API

  * TeamID: Integer 
  * Date:   String (YYYY-MM-DD)

  ```shell
      curl --location --request GET 'http://localhost:8080/api/v1/teams/{TeamID}/schedule?date={Date}'
  ```
  Example:
        ```shell
        curl --location --request GET 'http://localhost:8080/api/v1/teams/147/schedule?date=2022-07-21'
        ```
        
## Unit Tests
  go to front_service folder & run below make command to run unit tests & get code coverage
   
   ```shell
        make test
   ```
  output:
  
    go test -cover
    [MLB BallPark Segregation Service API]- unit testing 22:55:52 there is no games for our fav team on given day 0
    [MLB BallPark Segregation Service API]- unit testing 22:55:52 the indices of the fav games: [2]
    [MLB BallPark Segregation Service API]- unit testing 22:55:52 the indices of the fav games: [3 1]
    [MLB BallPark Segregation Service API]- unit testing 22:55:52 the indices of the fav games: [3 1]
    [MLB BallPark Segregation Service API]- unit testing 22:55:52 the indices of the fav games: [1 3]
    [MLB BallPark Segregation Service API]- unit testing 22:55:52 the indices of the fav games: [1 3]
    [MLB BallPark Segregation Service API]- unit testing 22:55:52 the indices of the fav games: [1 3]
    [MLB BallPark Segregation Service API]- unit testing 22:55:52 the indices of the fav games: [3 1]
    PASS
    coverage: 75.0% of statements
    ok  	github.com/mlb/mlb-ballpark-segregation-service/front_service	0.013s

## Development

if there are any changes to design file, run the below command in front_service folder to generate goa files
 ```shell
  make generate
 ```

Run the below command in /front-service folder to generate design files

```shell
goa gen github.com/mlb/mlb-ballpark-segregation-service/front_service/design
```
