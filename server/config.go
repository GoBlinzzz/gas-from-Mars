package server

import (
	"encoding/json"
	"net/http"
)

func GetConfig(w http.ResponseWriter, r *http.Request) Config {
	defer r.Body.Close()

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var config Config
	_ = decoder.Decode(&config)

	return config
}

func UpdateConfig(config *Config) {
	month := config.AmountOfClientsPerDay * config.TimeList.CarService * 30
	upd := config.TimeList.BuildingFillingPlace / month
	if config.TimeList.BuildingFillingPlace%month > 0 {
		upd++
	}
	config.TimeList.BuildingFillingPlace = upd
	upd = config.TimeList.BuildingStation / month
	if config.TimeList.BuildingStation%month > 0 {
		upd++
	}
	config.TimeList.BuildingStation = upd
}

func SetStatus(config Config) (status Status) {
	status = Status{
		FuelBalance:      config.StorageFuelBalance,
		MoneyBalance:     0,
		SuppliedStations: make([]bool, config.StationsAmount),
		WorkingTankers:   0,
		Stations:         setStations(config),
		Tankers:          setTankers(config),
	}
	return
}

func setStations(config Config) (stations []PetrolStation) {
	for i := 0; i < int(config.StationsAmount); i++ {
		stations = append(stations, PetrolStation{
			CashiersAmount: 1,
			FillingPlaces:  []int32{},
			FuelBalance:    config.StationFuelBalance,
		})
	}
	return
}

func setTankers(config Config) (tankers []Tanker) {
	for i := 0; i < int(config.TankersAmount); i++ {
		tankers = append(tankers, Tanker{
			Working:     false,
			Returning:   false,
			timeLeft:    0,
			FuelBalance: 0,
		})
	}
	return
}
