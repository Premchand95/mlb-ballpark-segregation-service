package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("Scheduler", func() {
	Description(`Scheduler service acts as a single source of truth for consumers to get a custom sorted schedule of games.`)
	HTTP(func() {
		Path("/api/v1/teams/{id}")
		Response("internal_error", StatusInternalServerError)
		Response("bad_gateway", StatusBadGateway)
		Response("bad_request", StatusBadRequest)
		Response("not_found", StatusNotFound)
	})
	Error("internal_error")
	Error("bad_gateway")
	Error("bad_request")
	Error("not_found")

	Method("index", func() {
		Meta("swagger:summary", `Retrieve a schedule of games`)
		Description(`Retrieves a schedule of games`)

		Payload(func() {
			Attribute("id", UInt, func() {
				Description(`The unique identifier of the team`)
				Example(120)
			})
			Attribute("date", String, "The date (YYYY-mm-dd) used to get all games scheduled on that day", func() {
				Format(FormatDate)
				Example("2022-08-10")
			})
			Required("id", "date")
		})

		Result(Schedule)
		HTTP(func() {
			GET("/schedule")
			Response(StatusOK)
		})
	})
})
