package mlbballparksegregationserviceapi

import (
	"context"
	"log"

	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
)

// Scheduler service example implementation.
// The example methods log the requests and return zero values.
type schedulersrvc struct {
	logger *log.Logger
}

// NewScheduler returns the Scheduler service implementation.
func NewScheduler(logger *log.Logger) scheduler.Service {
	return &schedulersrvc{logger}
}

// Retrieves a schedule of games
func (s *schedulersrvc) Index(ctx context.Context, p *scheduler.IndexPayload) (res *scheduler.Schedule, err error) {
	res = &scheduler.Schedule{}
	s.logger.Print("scheduler.index")
	return
}
