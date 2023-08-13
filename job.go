package vroom

type Job struct {
	//Id            int       `json:"id"`
	//Description   string    `json:"description"`
	//Location      *Location `json:"location"`
	//LocationIndex *int      `json:"location_index"`
	//locationID    int
	//Setup         int     `json:"setup"`
	//Service       int     `json:"service"`
	ShipmentStep
	TimeWindows [][]int `json:"time_windows,omitempty"`
	Delivery    []int   `json:"delivery,omitempty"`
	Pickup      []int   `json:"pickup,omitempty"`
	Skills      []int   `json:"skills,omitempty"`
	Priority    int     `json:"priority"`

	problem *Problem
}

func (j *Job) SetTimeWindows(tw [][]int) *Job { j.TimeWindows = tw; return j }
func (j *Job) SetDelivery(d []int) *Job       { j.Delivery = d; return j }
func (j *Job) SetPickup(p []int) *Job         { j.Pickup = p; return j }
func (j *Job) SetSkills(s []int) *Job         { j.Skills = s; return j }
func (j *Job) SetPriority(p int) *Job         { j.Priority = p; return j }

func (j *Job) SetLocationId(locationId int) *Job {
	j.locationID = &locationId
	j.problem.addLocationRef(locationId)
	return j
}

func (j *Job) SetMeasuredDelivery(m Measurement, val int) *Job {
	j.Delivery = j.problem.SetCapacity(m, val, j.Delivery)
	return j
}

func (j *Job) SetMeasuredPickup(m Measurement, val int) *Job {
	j.Pickup = j.problem.SetCapacity(m, val, j.Pickup)
	return j
}
