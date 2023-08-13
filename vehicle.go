package vroom

type VehicleProfile string

const (
	VPCar VehicleProfile = "car"
)

type Vehicle struct {
	Id          int            `json:"id"`
	Profile     VehicleProfile `json:"profile"`
	Description string         `json:"description"`
	Start       *Location      `json:"start"`
	StartIndex  *int           `json:"start_index"`
	startID     *int
	End         *Location `json:"end"`
	EndIndex    *int      `json:"end_index"`
	endID       *int
	Capacity    []int `json:"capacity"`
	Costs       []int `json:"costs"`
	Skills      []int `json:"skills"`
	TimeWindow  []int `json:"time_window"`
	Breaks      []int `json:"breaks"`
	problem     *Problem
}

func (v *Vehicle) SetProfile(p VehicleProfile) *Vehicle { v.Profile = p; return v }
func (v *Vehicle) SetDescription(d string) *Vehicle     { v.Description = d; return v }
func (v *Vehicle) SetStart(l *Location) *Vehicle        { v.Start = l; return v }
func (v *Vehicle) SetStartIndex(idx int) *Vehicle       { v.StartIndex = &idx; return v }
func (v *Vehicle) SetEnd(l *Location) *Vehicle          { v.End = l; return v }
func (v *Vehicle) SetEndIndex(idx int) *Vehicle         { v.EndIndex = &idx; return v }
func (v *Vehicle) SetCapacity(c []int) *Vehicle         { v.Capacity = c; return v }
func (v *Vehicle) SetCosts(c []int) *Vehicle            { v.Costs = c; return v }
func (v *Vehicle) SetSkills(s []int) *Vehicle           { v.Skills = s; return v }
func (v *Vehicle) SetTimeWindow(tw []int) *Vehicle      { v.TimeWindow = tw; return v }
func (v *Vehicle) SetBreaks(b []int) *Vehicle           { v.Breaks = b; return v }

func (v *Vehicle) SetStartId(locationId int) *Vehicle {
	v.startID = &locationId
	v.problem.addLocationRef(locationId)
	return v
}

func (v *Vehicle) SetEndId(locationId int) *Vehicle {
	v.endID = &locationId
	v.problem.addLocationRef(locationId)
	return v
}

func (v *Vehicle) SetMeasuredCapacity(m Measurement, val int) *Vehicle {
	v.Capacity = v.problem.SetCapacity(m, val, v.Capacity)
	return v
}
