// Code generated by goa v3.8.2, DO NOT EDIT.
//
// Scheduler endpoints
//
// Command:
// $ goa gen
// github.com/mlb/mlb-ballpark-segregation-service/front_service/design

package scheduler

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "Scheduler" service endpoints.
type Endpoints struct {
	Index goa.Endpoint
}

// NewEndpoints wraps the methods of the "Scheduler" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Index: NewIndexEndpoint(s),
	}
}

// Use applies the given middleware to all the "Scheduler" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Index = m(e.Index)
}

// NewIndexEndpoint returns an endpoint function that calls the method "index"
// of service "Scheduler".
func NewIndexEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*IndexPayload)
		res, err := s.Index(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedSchedule(res, "default")
		return vres, nil
	}
}