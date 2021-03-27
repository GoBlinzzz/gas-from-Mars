package server

import "fmt"

func Pay(status *Status, log *[]string) {
	for i, val := range status.Stations {
		status.MoneyBalance -= val.MaintenancePrice
		*log = append(*log, fmt.Sprintf("Spent %d for maintaining %d gas station", val.MaintenancePrice, i))
	}
}

func CountMaintenancePrice(status *Status, config Config) {
	for i, val := range status.Stations {
		status.Stations[i].MaintenancePrice = config.PriceList.StationMaintaining
		status.Stations[i].MaintenancePrice += config.PriceList.FillingPlaceMaintaining * int32(len(val.FillingPlaces))
		status.Stations[i].MaintenancePrice += config.Salaries.Manager
		status.Stations[i].MaintenancePrice += config.Salaries.Guard
		status.Stations[i].MaintenancePrice += config.Salaries.Cashier * val.CashiersAmount
		status.Stations[i].MaintenancePrice += config.Salaries.PumpAttendant * int32(len(val.FillingPlaces))
	}
}

func CountOrdinaryProfit(profit int32, fillingPlaces int, rate float32) int32 {
	if fillingPlaces == 0 {
		return int32(float32(profit) * 0.3)
	} else {
		return profit * int32(1+float32(fillingPlaces)*rate)
	}
}

func CopyStatus(status *Status, statusCopy *Status) {
	statusCopy.FuelBalance = status.FuelBalance
	statusCopy.MoneyBalance = status.MoneyBalance
	for _, val := range status.NewStations {
		statusCopy.NewStations = append(statusCopy.NewStations, val)
	}
	for _, val := range status.NewFillingPlaces {
		statusCopy.NewFillingPlaces = append(statusCopy.NewFillingPlaces, val)
	}
	for _, val := range status.Tankers {
		statusCopy.Tankers = append(statusCopy.Tankers, val)
	}
	for _, val := range status.Stations {
		statusCopy.Stations = append(statusCopy.Stations, val)
	}
}

func ChooseStrategy(status *Status, config Config, currentMonth int) {
	var statusCopy Status
	CopyStatus(status, &statusCopy)
	for i := range status.Stations {
		tankersTripsNumber := len(status.Tankers) * int(config.AmountOfClientsPerDay*config.TimeList.CarService*30) / 2 / int(config.TimeList.GasDelivery)
		if currentMonth+int(config.TimeList.BuildingFillingPlace) < 11 {
			status.NewFillingPlaces[currentMonth+int(config.TimeList.BuildingFillingPlace)][i] = true
			if ProfitPrediction(statusCopy, config, currentMonth, true) < ProfitPrediction(statusCopy, config, currentMonth, false) {
				status.NewFillingPlaces[currentMonth+int(config.TimeList.BuildingFillingPlace)][i] = false
			} else if tankersTripsNumber < len(status.Stations) {
				status.Tankers = append(status.Tankers, Tanker{
					Working:        false,
					Returning:      false,
					timeLeft:       0,
					DestinationNum: 0,
					Destination:    nil,
					FuelBalance:    0,
				})
			}
		}
	}
	ind := true
	tankersAmount := len(status.Tankers)
	for ind {
		tankersTripsNumber := len(status.Tankers) * int(config.AmountOfClientsPerDay*config.TimeList.CarService*30) / 2 / int(config.TimeList.GasDelivery)
		if currentMonth+1 < 12 {
			status.NewStations[currentMonth+int(config.TimeList.BuildingStation)]++
			if status.NewStations[currentMonth+int(config.TimeList.BuildingStation)] < tankersTripsNumber {
				if ProfitPrediction(statusCopy, config, currentMonth, true) < ProfitPrediction(statusCopy, config, currentMonth, false) {
					status.NewStations[currentMonth+int(config.TimeList.BuildingStation)]--
					ind = !ind
				}
			} else {
				if ProfitPrediction(statusCopy, config, currentMonth, true) < ProfitPrediction(statusCopy, config, currentMonth, false)+config.PriceList.TankerPrice*int32(len(status.Tankers)-tankersAmount+1) {
					status.NewStations[currentMonth+int(config.TimeList.BuildingStation)]--
					ind = !ind
				} else if tankersTripsNumber > status.NewStations[currentMonth+int(config.TimeList.BuildingStation)]+len(status.Stations) {
					status.Tankers = append(status.Tankers, Tanker{
						Working:        false,
						Returning:      false,
						timeLeft:       0,
						DestinationNum: 0,
						Destination:    nil,
						FuelBalance:    0,
					})
				} else {
					ind = !ind
				}
			}
		} else {
			ind = !ind
		}
	}
}

func ProfitPrediction(status Status, config Config, currentMonth int, flag bool) int32 {
	prediction := int32(0)
	fuelBalance := status.FuelBalance
	if flag {
		prediction -= FinishBuildingNewObjects(&status, config, currentMonth+1)
	}
	for j, val := range status.Stations {
		maxFuel := config.AmountOfClientsPerDay * config.FuelAmountPerClient * 30 * int32(len(val.FillingPlaces))
		if maxFuel > val.FuelBalance {
			prediction += val.FuelBalance * CountOrdinaryProfit(config.PriceList.OrdinaryProfit, len(val.FillingPlaces), config.FillingPlacesRatio)
			maxFuel -= val.FuelBalance
			status.Stations[j].FuelBalance = 0
			if fuelBalance > maxFuel {
				prediction += maxFuel * CountOrdinaryProfit(config.PriceList.OrdinaryProfit, len(val.FillingPlaces), config.FillingPlacesRatio)
				fuelBalance -= maxFuel
			} else {
				prediction += fuelBalance * CountOrdinaryProfit(config.PriceList.OrdinaryProfit, len(val.FillingPlaces), config.FillingPlacesRatio)
				fuelBalance = 0
			}
		} else {
			prediction += maxFuel * CountOrdinaryProfit(config.PriceList.OrdinaryProfit, len(val.FillingPlaces), config.FillingPlacesRatio)
			status.Stations[j].FuelBalance -= maxFuel
		}
	}
	return prediction
}
