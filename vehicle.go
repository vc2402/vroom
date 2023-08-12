package vroom

type VehicleProfile string

const (
	VPCar VehicleProfile = "car"
)

type Vehicle struct {
	Id         int            `json:"id"`
	Profile    VehicleProfile `json:"profile"`
	Start      []float64      `json:"start"`
	StartIndex int            `json:"start_index"`
	End        []float64      `json:"end"`
	EndIndex   int            `json:"end_index"`
	Capacity   []int          `json:"capacity"`
	Costs      []int          `json:"costs"`
	Skills     []int          `json:"skills"`
	TimeWindow []int          `json:"time_window"`
	Breaks     []int          `json:"breaks"`
}

type Vehicles []*Vehicle

func (vc *Vehicles) NewVehicle() *Vehicle {

}
