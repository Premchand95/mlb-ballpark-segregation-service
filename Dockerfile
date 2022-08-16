FROM golang:1.18.2-alpine as build
WORKDIR /go/src/github.com/mlb/mlb-ballpark-segregation-service
COPY . .
RUN go mod download
RUN go build -o front_service_binary /go/src/github.com/mlb/mlb-ballpark-segregation-service/front_service/cmd/front_service

FROM alpine:latest 
COPY --from=build /go/src/github.com/mlb/mlb-ballpark-segregation-service /mlb
EXPOSE 8080
CMD [ "./mlb/front_service_binary"]



