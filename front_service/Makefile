# to generate gen files
generate:
	goa gen github.com/mlb/mlb-ballpark-segregation-service/front_service/design
# to run unit tests
test:
	go test -cover
# to create mock clients for services
mocks::
	mockgen -source services/statsAPI/client.go -destination mocks/mlb_statsAPI/mlb_statsAPI_mock_client.go -package mock_statsapi 
	mockgen -source services/requests/client.go -destination mocks/mlb_requests/mlb_requests_mock_client.go -package mock_requests