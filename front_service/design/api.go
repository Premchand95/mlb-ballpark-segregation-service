//Copyright 2022 MLB Advanced Media, L.P.  Use of any content on this page acknowledges agreement to the terms posted here http://gdx.mlb.com/components/copyright.txt
package design

import (
	"fmt"

	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/docs"
)

const (
	ServiceName = "MLB BallPark Segregation Service API"
	// MajorVersion is incremented for breaking changes
	MajorVersion = 1
	// MinorVersion is incremented for new, backwards-compatible functionality is added
	MinorVersion = 0
	// PatchVersion is incremented for major bug fixes
	PatchVersion = 0
)

var _ = API(ServiceName, func() {
	Title("MLB BallPark Segregation Service API")
	Version(fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion))
	Description(`The MLB BallPark Segregation Service API provides REST-ful API to retrieve custom schedules of mlb games`)
})
