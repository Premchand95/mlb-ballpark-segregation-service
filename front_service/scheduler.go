package mlbballparksegregationserviceapi

import (
	"context"
	"log"

	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
	stats "github.com/mlb/mlb-ballpark-segregation-service/front_service/services/statsAPI"
)

// Scheduler service example implementation.
// The example methods log the requests and return zero values.
type schedulersrvc struct {
	logger *log.Logger
}

// checks weather all endpoints are implemented or not
var (
	_ scheduler.Service = (*schedulersrvc)(nil)
)

// NewScheduler returns the Scheduler service implementation.
func NewScheduler(logger *log.Logger) scheduler.Service {
	return &schedulersrvc{logger}
}

// Retrieves a schedule of games
func (s *schedulersrvc) Index(ctx context.Context, p *scheduler.IndexPayload) (res *scheduler.Schedule, err error) {
	res = &scheduler.Schedule{}

	statsC, err := stats.NewStatsClient(ctx)
	if err != nil {
		s.logger.Print("error", err)
		return nil, err
	}

	res, err = statsC.GetStatsAPISchedule(ctx, p.Date)
	if err != nil {
		s.logger.Print("error", err)
		return nil, err
	}

	s.logger.Print("scheduler.index")
	return
}
