package vroom

type ShipmentStep struct {
	Id            int       `json:"id"`
	Description   string    `json:"description"`
	Location      *Location `json:"location"`
	LocationIndex *int      `json:"location_index"`
	locationID    *int
	Setup         int     `json:"setup"`
	Service       int     `json:"service"`
	TimeWindows   [][]int `json:"time_windows"`
}

type Shipment struct {
	Pickup   ShipmentStep `json:"pickup"`
	Delivery ShipmentStep `json:"delivery"`
	Amount   []int        `json:"amount"`
	Skills   []int        `json:"skills"`
	Priority int          `json:"priority"`
	problem  *Problem
}

func (ss *ShipmentStep) SetDescription(d string) *ShipmentStep   { ss.Description = d; return ss }
func (ss *ShipmentStep) SetLocation(l *Location) *ShipmentStep   { ss.Location = l; return ss }
func (ss *ShipmentStep) SetLocationIndex(idx int) *ShipmentStep  { ss.LocationIndex = &idx; return ss }
func (ss *ShipmentStep) SetSetup(s int) *ShipmentStep            { ss.Setup = s; return ss }
func (ss *ShipmentStep) SetService(s int) *ShipmentStep          { ss.Service = s; return ss }
func (ss *ShipmentStep) SetTimeWindows(tw [][]int) *ShipmentStep { ss.TimeWindows = tw; return ss }

func (s *Shipment) SetAmount(a []int) *Shipment      { s.Amount = a; return s }
func (s *Shipment) SetSkills(skills []int) *Shipment { s.Skills = skills; return s }
func (s *Shipment) SetPriority(p int) *Shipment      { s.Priority = p; return s }

func (s *Shipment) SetPickupLocationId(locationId int) *Shipment {
	s.Pickup.locationID = &locationId
	s.problem.addLocationRef(locationId)
	return s
}

func (s *Shipment) SetDeliveryLocationId(locationId int) *Shipment {
	s.Delivery.locationID = &locationId
	s.problem.addLocationRef(locationId)
	return s
}

func (s *Shipment) SetMeasuredAmount(m Measurement, val int) *Shipment {
	s.Amount = s.problem.SetCapacity(m, val, s.Amount)
	return s
}
