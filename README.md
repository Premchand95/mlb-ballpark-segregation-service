# MLB-Ballpark-services

The MLB BallPark Segregation Service API provides REST-ful API to retrieve custom schedules of mlb games

## Prerequisites

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

## Goa Generate 

Run the below command in /front-service folder to generate design files

```shell
goa gen github.com/mlb/mlb-ballpark-segregation-service/front_service/design
```
