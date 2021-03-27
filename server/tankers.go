package server

import (
	"sync"
)

func ManageTankers(status *Status, config Config, currentMonth int) {
	var localWG sync.WaitGroup

	for i := range status.Tankers {
		localWG.Add(1)
		if status.Tankers[i].Working {
			go ProcessingTanker(&status.Tankers[i], status, config, &localWG)
		} else if currentMonth != 11 {
			go spreadFuel(status, config, currentMonth, &localWG)
		} else {
			localWG.Done()
		}
	}
	localWG.Wait()
}

func spreadFuel(status *Status, config Config, currentMonth int, wg *sync.WaitGroup) {
	var (
		freeTankers    []int
		chosenStations []int
	)
	for i := 0; i < len(status.Tankers); i++ {
		if !status.Tankers[i].Working {
			freeTankers = append(freeTankers, i)
		}
	}
	for i := 0; len(chosenStations) != len(freeTankers); i++ {
		choice := chooseStation(status, config)
		chosenStations = append(chosenStations, choice)
		status.SuppliedStations[choice] = true
	}
	for i, t := range freeTankers {
		launchTanker(status, config, &status.Tankers[t], chosenStations[i], currentMonth)
	}
	wg.Done()
}

func launchTanker(status *Status, config Config, tanker *Tanker, destination int, currentMonth int) {
	tanker.Working = true
	tanker.Returning = false
	tanker.DestinationNum = destination
	tanker.Destination = &status.Stations[destination]
	tanker.timeLeft = config.TimeList.GasDelivery
	neededFuel := predictFuelNeed(status, config, destination, currentMonth)
	if neededFuel == 0 {
		tanker.Working = false
	} else {
		if neededFuel > status.FuelBalance {
			tanker.FuelBalance = status.FuelBalance
			status.FuelBalance = 0
		} else {
			tanker.FuelBalance = neededFuel
			status.FuelBalance -= neededFuel
		}
	}
}

func predictFuelNeed(status *Status, config Config, destination int, currentMonth int) int32 {
	currentBalance := status.Stations[destination].FuelBalance
	capacity := config.FuelAmountPerClient * config.AmountOfClientsPerDay * 30
	if status.NewFillingPlaces[currentMonth+1][destination] {
		capacity *= int32(len(status.Stations[destination].FillingPlaces) + 1)
	} else {
		capacity *= int32(len(status.Stations[destination].FillingPlaces))
	}
	capacity -= currentBalance
	if capacity > 0 {
		return capacity
	}
	return 0
}

func ProcessingTanker(tanker *Tanker, status *Status, config Config, wg *sync.WaitGroup) {
	if tanker.Returning {
		if tanker.timeLeft == 1 {
			tanker.Working = false
			status.WorkingTankers--
		}
		tanker.timeLeft--
	} else {
		if tanker.timeLeft == 1 {
			tanker.Destination.FuelBalance += tanker.FuelBalance
			tanker.Returning = true
			tanker.timeLeft = config.TimeList.GasDelivery
			status.SuppliedStations[tanker.DestinationNum] = false
		} else {
			tanker.timeLeft--
		}
	}
	wg.Done()
}

func predictStationFuelConsumption(st PetrolStation, config Config) int32 {
	return config.TimeList.GasDelivery * 3 / config.TimeList.CarService * int32(len(st.FillingPlaces)) * config.FuelAmountPerClient
}

func chooseStation(status *Status, config Config) int {
	choice := 0
	fillingPlaces := 0

	for i, st := range status.Stations {
		if !status.SuppliedStations[i] && len(st.FillingPlaces) > fillingPlaces && 0 >= st.FuelBalance-predictStationFuelConsumption(st, config) {
			choice = i
			fillingPlaces = len(st.FillingPlaces)
		}
	}
	return choice
}
