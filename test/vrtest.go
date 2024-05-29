package main

import (
	"fmt"
	"github.com/vc2402/vroom"
)

func main() {
	client := vroom.NewClient("http://localhost:8000/vroom")
	locations := map[int][2]float64{
		10: {25.28814777028492, 54.65403824495201},
		12: {25.59632682800293, 55.210960388183594},
		14: {25.60372543334961, 55.93915557861328},
	}
	problem := &vroom.Problem{}
	//timeWindow := []int{0, 7}
	problem.NewVehicle().SetStartId(10).SetMeasuredCapacity(vroom.MeasWeight, 1000).SetMeasuredCapacity(vroom.MeasVolume, 100) //.SetTimeWindow(timeWindow)
	problem.NewVehicle().SetStartId(10)                                                                                        //.SetTimeWindow(timeWindow)
	//shipment :=
	problem.NewShipment().SetPickupLocationId(10).SetDeliveryLocationId(12).
		SetMeasuredAmount(vroom.MeasWeight, 200).SetMeasuredAmount(vroom.MeasVolume, 20)
	//shipment.Pickup.SetTimeWindows([][]int{timeWindow})
	//shipment.Delivery.SetTimeWindows([][]int{timeWindow})

	//shipment =
	problem.NewShipment().SetPickupLocationId(10).SetDeliveryLocationId(14).
		SetMeasuredAmount(vroom.MeasWeight, 200).SetMeasuredAmount(vroom.MeasVolume, 20)
	//shipment.Pickup.SetTimeWindows([][]int{timeWindow})
	//shipment.Delivery.SetTimeWindows([][]int{timeWindow})
	//problem.NewJob().SetLocationId(10).SetMeasuredPickup(vroom.MeasWeight, 1000).SetMeasuredPickup(vroom.MeasVolume, 100)
	//problem.NewJob().SetLocationId(12).SetMeasuredDelivery(vroom.MeasWeight, 300).SetMeasuredDelivery(vroom.MeasVolume, 30)
	//problem.NewJob().SetLocationId(12).SetMeasuredDelivery(vroom.MeasWeight, 400).SetMeasuredDelivery(vroom.MeasVolume, 40)
	//problem.NewJob().SetLocationId(14).SetMeasuredDelivery(vroom.MeasWeight, 300).SetMeasuredDelivery(vroom.MeasVolume, 30)

	problem.FillLocations(func(locationID int) vroom.Location {
		return locations[locationID]
	})

	sol, err := client.Solve(problem)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("Solution: %+v/n", *sol)
}
