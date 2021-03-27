package server

import (
	"sync"
	"sync/atomic"
)

func ProcessStations(status *Status, config Config) {
	var localWG sync.WaitGroup
	for i := range status.Stations {
		localWG.Add(1)
		go ProcessingStation(&status.Stations[i], &status.MoneyBalance, config, &localWG)
	}
	localWG.Wait()
}

func ProcessingStation(st *PetrolStation, bank *int32, config Config, wg *sync.WaitGroup) {
	for i, fpl := range st.FillingPlaces {
		if fpl == 1 {
			atomic.AddInt32(bank, CountOrdinaryProfit(config.PriceList.OrdinaryProfit, len(st.FillingPlaces), config.FillingPlacesRatio))
			st.FillingPlaces[i]--
		} else if fpl == 0 && st.FuelBalance >= config.FuelAmountPerClient {
			atomic.AddInt32(&st.FuelBalance, -config.FuelAmountPerClient)
			st.FillingPlaces[i] += config.TimeList.CarService
		} else if fpl != 0 {
			st.FillingPlaces[i]--
		}
	}
	wg.Done()
}

func FinishBuildingNewObjects(status *Status, config Config, currentMonth int) int32 {
	price := int32(0)
	for i := 0; i < status.NewStations[currentMonth]; i++ {
		status.Stations = append(status.Stations, PetrolStation{
			CashiersAmount: 1,
			FillingPlaces:  nil,
			FuelBalance:    0,
		})
		status.SuppliedStations = append(status.SuppliedStations, false)
		price += config.PriceList.StationMaintaining
		price += config.Salaries.Manager
		price += config.Salaries.Guard
		price += config.Salaries.Cashier
	}
	for key, val := range status.NewFillingPlaces[currentMonth] {
		if val {
			status.MoneyBalance -= config.PriceList.FillingPlacePrice
			status.Stations[key].FillingPlaces = append(status.Stations[key].FillingPlaces, 0)
			price += config.Salaries.PumpAttendant
			if len(status.Stations[key].FillingPlaces) > int(status.Stations[key].CashiersAmount*config.CashierCapacity) {
				status.Stations[key].CashiersAmount++
				price += config.Salaries.Cashier
			}
		}
	}
	return price
}
