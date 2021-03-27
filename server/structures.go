package server

import "time"

type Config struct {
	Deliveries            []int32
	StorageFuelBalance    int32
	StationFuelBalance    int32
	StationsAmount        int32
	TankersAmount         int32
	AmountOfClientsPerDay int32
	FuelAmountPerClient   int32
	PriceList             PriceList
	TimeList              TimeList
	Salaries              Salaries
	FillingPlacesRatio    float32
	CashierCapacity       int32
	SackProbability       float32
	EquivalentOfMonth     time.Duration
}

type PetrolStation struct {
	CashiersAmount   int32
	FillingPlaces    []int32
	FuelBalance      int32
	MaintenancePrice int32
}

type Salaries struct {
	Manager       int32
	Cashier       int32
	Guard         int32
	PumpAttendant int32
}

type TimeList struct {
	GasDelivery          int32
	CarService           int32
	BuildingStation      int32
	BuildingFillingPlace int32
}

type PriceList struct {
	StationMaintaining      int32
	FillingPlaceMaintaining int32
	TankerPrice             int32
	FillingPlacePrice       int32
	OrdinaryProfit          int32
}
type Tanker struct {
	Working        bool
	Returning      bool
	timeLeft       int32
	DestinationNum int
	Destination    *PetrolStation
	FuelBalance    int32
}

type Status struct {
	FuelBalance      int32
	MoneyBalance     int32
	SuppliedStations []bool
	WorkingTankers   int
	NewStations      []int
	NewFillingPlaces []map[int]bool
	Tankers          []Tanker
	Stations         []PetrolStation
}
