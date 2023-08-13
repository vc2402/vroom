package vroom

type Summary struct {
	Cost        int           `json:"cost"`
	Routes      int           `json:"routes"`
	Unassigned  int           `json:"unassigned"`
	Delivery    []int         `json:"delivery"`
	Pickup      []int         `json:"pickup"`
	Setup       int           `json:"setup"`
	Service     int           `json:"service"`
	Duration    int           `json:"duration"`
	WaitingTime int           `json:"waiting_time"`
	Priority    int           `json:"priority"`
	Distance    int           `json:"distance"`
	Violations  []interface{} `json:"violations"`
}

type RouteStep struct {
	Type        string        `json:"type"`
	Location    []float64     `json:"location,omitempty"`
	Setup       int           `json:"setup"`
	Service     int           `json:"service"`
	WaitingTime int           `json:"waiting_time"`
	Load        []int         `json:"load"`
	Arrival     int           `json:"arrival"`
	Duration    int           `json:"duration"`
	Violations  []interface{} `json:"violations"`
	Distance    int           `json:"distance"`
	Id          int           `json:"id,omitempty"`
	Job         int           `json:"job,omitempty"`
}

type Route struct {
	Vehicle     int           `json:"vehicle"`
	Cost        int           `json:"cost"`
	Delivery    []int         `json:"delivery"`
	Pickup      []int         `json:"pickup"`
	Setup       int           `json:"setup"`
	Service     int           `json:"service"`
	Duration    int           `json:"duration"`
	WaitingTime int           `json:"waiting_time"`
	Priority    int           `json:"priority"`
	Distance    int           `json:"distance"`
	Steps       []RouteStep   `json:"steps"`
	Violations  []interface{} `json:"violations"`
	Geometry    string        `json:"geometry"`
}

type Solution struct {
	Code       int `json:"code"`
	Summary    `json:"summary"`
	Unassigned []interface{} `json:"unassigned"`
	Routes     []Route       `json:"routes"`
}
