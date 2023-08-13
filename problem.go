package vroom

import "errors"

type Location [2]float64

type Measurement int

const (
	MeasVolume Measurement = iota
	MeasWeight
	MeasAre
	MeasPLL
)

var ErrLocationNotDefined = errors.New("location is not defined")

type Problem struct {
	Vehicles     []*Vehicle  `json:"vehicles"`
	Jobs         []*Job      `json:"jobs,omitempty"`
	Shipments    []*Shipment `json:"shipments,omitempty"`
	locations    map[int]Location
	measurements map[Measurement]int
	jobsCount    int
	err          error
}

func (p *Problem) NewVehicle(id ...int) *Vehicle {
	var vehId int
	if len(id) > 0 {
		vehId = id[0]
	} else {
		vehId = len(p.Vehicles) + 1
	}
	v := &Vehicle{Id: vehId, problem: p}
	p.AddVehicle(v)
	return v
}

func (p *Problem) AddVehicle(v *Vehicle) *Problem {
	p.Vehicles = append(p.Vehicles, v)
	return p
}

func (p *Problem) NewJob(id ...int) *Job {
	var jid int
	if len(id) > 0 {
		jid = id[0]
	} else {
		jid = p.NextJobId()
	}
	j := &Job{ShipmentStep: ShipmentStep{Id: jid}, problem: p}
	p.AddJob(j)
	return j
}

func (p *Problem) AddJob(j *Job) *Problem {
	p.Jobs = append(p.Jobs, j)
	return p
}

func (p *Problem) NewShipment(id ...int) *Shipment {
	var jid1 int
	if len(id) > 0 {
		jid1 = id[0]
	} else {
		jid1 = p.NextJobId()
	}
	var jid2 int
	if len(id) > 1 {
		jid2 = id[1]
	} else {
		jid2 = p.NextJobId()
	}
	j := &Shipment{Pickup: ShipmentStep{Id: jid1}, Delivery: ShipmentStep{Id: jid2}, problem: p}
	p.AddShipment(j)
	return j
}

func (p *Problem) AddShipment(s *Shipment) *Problem {
	p.Shipments = append(p.Shipments, s)
	return p
}

func (p *Problem) GetMeasurementIndex(measurement Measurement) int {
	if idx, ok := p.measurements[measurement]; ok {
		return idx
	}
	idx := len(p.measurements)
	if idx == 0 {
		p.measurements = map[Measurement]int{measurement: idx}
	} else {
		p.measurements[measurement] = idx
	}
	return idx
}

func (p *Problem) SetCapacity(measurement Measurement, val int, capacities []int) []int {
	idx := p.GetMeasurementIndex(measurement)
	if len(capacities) <= idx {
		newSlice := make([]int, len(p.measurements))
		copy(newSlice, capacities)
		capacities = newSlice
	}
	capacities[idx] = val
	return capacities
}

func (p *Problem) NextJobId() (id int) {
	id = p.jobsCount + 1
	p.jobsCount++
	return
}

func (p *Problem) FillLocations(locationResolver func(locationID int) Location) bool {
	for id := range p.locations {
		loc := locationResolver(id)
		if loc[0] == 0 && loc[1] == 0 {
			return false
		}
		p.locations[id] = loc
	}
	return true
}

func (p *Problem) GetLocation(id int) (loc Location, err error) {
	if p.locations == nil {
		err = ErrLocationNotDefined
		return
	}
	loc, ok := p.locations[id]
	if !ok {
		err = ErrLocationNotDefined
	}
	return
}

func (p *Problem) Err() error { return p.err }

func (p *Problem) processLocations() {
	if len(p.locations) > 0 {
		for _, vehicle := range p.Vehicles {
			if vehicle.startID != nil {
				var loc Location
				loc, p.err = p.GetLocation(*vehicle.startID)
				vehicle.Start = &loc
				if p.err != nil {
					return
				}
			}
			if vehicle.endID != nil {
				var loc Location
				loc, p.err = p.GetLocation(*vehicle.endID)
				vehicle.End = &loc
				if p.err != nil {
					return
				}
			}
		}
		for _, job := range p.Jobs {
			p.processShipmentStepLocation(&job.ShipmentStep)
			if p.err != nil {
				return
			}
		}
		for _, shipment := range p.Shipments {
			p.processShipmentStepLocation(&shipment.Delivery)
			if p.err != nil {
				return
			}
			p.processShipmentStepLocation(&shipment.Pickup)
			if p.err != nil {
				return
			}
		}
	}
}

func (p *Problem) processShipmentStepLocation(ss *ShipmentStep) {
	if ss.locationID != nil {
		var loc Location
		loc, p.err = p.GetLocation(*ss.locationID)
		ss.Location = &loc
	}
}

func (p *Problem) addLocationRef(locationId int) {
	if p.locations == nil {
		p.locations = map[int]Location{}
	}
	if _, ok := p.locations[locationId]; !ok {
		p.locations[locationId] = Location{}
	}
}
