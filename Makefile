#to build docker image
docker-image:
	docker build -t mlb-segregation-service .
# to build binary
build:
	go build -o front_service_binary front_service/cmd/front_service/*
# to run the application
run:
	docker run -p 8080:8080 mlb-segregation-service:latest


