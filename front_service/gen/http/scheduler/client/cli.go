// Code generated by goa v3.8.2, DO NOT EDIT.
//
// Scheduler HTTP client CLI support package
//
// Command:
// $ goa gen
// github.com/mlb/mlb-ballpark-segregation-service/front_service/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
	goa "goa.design/goa/v3/pkg"
)

// BuildIndexPayload builds the payload for the Scheduler index endpoint from
// CLI flags.
func BuildIndexPayload(schedulerIndexBody string, schedulerIndexID string) (*scheduler.IndexPayload, error) {
	var err error
	var body IndexRequestBody
	{
		err = json.Unmarshal([]byte(schedulerIndexBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"date\": \"2022-08-10\"\n   }'")
		}
		err = goa.MergeErrors(err, goa.ValidateFormat("body.date", body.Date, goa.FormatDate))

		if err != nil {
			return nil, err
		}
	}
	var id uint
	{
		var v uint64
		v, err = strconv.ParseUint(schedulerIndexID, 10, strconv.IntSize)
		id = uint(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be UINT")
		}
	}
	v := &scheduler.IndexPayload{
		Date: body.Date,
	}
	v.ID = id

	return v, nil
}