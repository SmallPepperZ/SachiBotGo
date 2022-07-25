package api

import "time"

type GlobalType struct {
	AppID string
	StartTime time.Time
}
var Globals = &GlobalType{}