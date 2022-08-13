package statsapi

import (
	"context"
	"encoding/json"
	"errors"
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

var _ Client = (*statsClient)(nil)

type (
	//Stats client interface
	Client interface {
		GetStatsAPISchedule(context.Context, string) (*types.Schedule, error)
	}

	// Stats API client struct
	statsClient struct {
		req    requests.Client
		logger log.Logger
	}

	StatsClientParams struct {
		Req    *requests.Client
		Logger *log.Logger
	}
)

//NewStatsClient initialize new Stats client
func NewStatsClient(params *StatsClientParams) (*statsClient, error) {
	if params.Req == nil || params.Logger == nil {
		return nil, errors.New("all parameters are required for instantiating the stats API client")
	}
	return &statsClient{
		req:    *params.Req,
		logger: *params.Logger,
	}, nil
}

func (s *statsClient) GetStatsAPISchedule(ctx context.Context, date string) (*types.Schedule, error) {
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
