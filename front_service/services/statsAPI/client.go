package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	types "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
	"github.com/mlb/mlb-ballpark-segregation-service/front_service/services/requests"
)

const (
	host     string = "https://statsapi.mlb.com/api/v1/schedule"
	sportId  string = "1"
	language string = "en"
)

var _ ClientStatsAPI = (*StatsClient)(nil)

type (
	// Stats client interface
	ClientStatsAPI interface {
		GetStatsAPISchedule(context.Context, string) (*types.Schedule, error)
	}

	// Stats API client struct
	StatsClient struct {
		req *requests.Client
	}
)

type (
	Client interface {
		GetStatsAPISchedule(context.Context, string) (*types.Schedule, error)
	}
)

//NewStatsClient initialize new Stats client
func NewStatsClient(ctx context.Context) (*StatsClient, error) {
	// create new requests client
	req, err := requests.NewClient()
	if err != nil {
		log.Println(ctx, "error", "error while creating new requests client", err.Error())
		return nil, fmt.Errorf("error while creating new requests client")
	}
	return &StatsClient{req: req}, err
}

func (s *StatsClient) GetStatsAPISchedule(ctx context.Context, date string) (*types.Schedule, error) {
	var schedule types.Schedule
	params := make(map[string]string, 0)

	params["date"] = date
	params["sportId"] = sportId
	params["language"] = language

	// prepare parameters
	res, err := s.req.Get(ctx, host, params)
	if err != nil {
		log.Println(ctx, "error", "error calling Get method", err.Error())
		return nil, fmt.Errorf("error calling GET method")
	}
	// unmarshal the GET Response
	err = json.Unmarshal(res, &schedule)
	if err != nil {
		log.Println(ctx, "error", "error while un marshalling GET response", err.Error())
		return nil, fmt.Errorf("error un marshalling Get response into type scheduler")
	}
	return &schedule, nil
}
