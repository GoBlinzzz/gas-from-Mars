package server

import (
	"time"
)

func GetConfig() (config Config) {
	config = Config{
		Deliveries:            []int32{1000, 1000, 2500, 1400, 5000, 3000, 200, 800, 450, 1000, 250, 450},
		StorageFuelBalance:    1000,
		StationFuelBalance:    500,
		StationsAmount:        2,
		TankersAmount:         0,
		AmountOfClientsPerDay: 20,
		FuelAmountPerClient:   50,
		PriceList: PriceList{
			StationMaintaining:      1000,
			FillingPlaceMaintaining: 70,
			TankerPrice:             200,
			FillingPlacePrice:       800,
			OrdinaryProfit:          150,
		},
		TimeList: TimeList{
			GasDelivery:          300,
			CarService:           1,
			BuildingStation:      20,
			BuildingFillingPlace: 6,
		},
		Salaries: Salaries{
			Manager:       2500,
			Cashier:       2000,
			Guard:         1700,
			PumpAttendant: 1300,
		},
		FillingPlacesRatio: 2,
		CashierCapacity:    3,
		EquivalentOfMonth:  time.Second * 10,
	}
	return
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
