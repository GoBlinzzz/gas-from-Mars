package main

import (
	"backend/server"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	router := mux.NewRouter()
	config := server.GetConfig()
	status := server.SetStatus(config)
	server.UpdateConfig(&config)
	var logger []string
	month := int(config.AmountOfClientsPerDay * config.TimeList.CarService * 30)
	status.NewStations = make([]int, 12)
	status.NewFillingPlaces = make([]map[int]bool, 12)
	for i := range status.NewFillingPlaces {
		status.NewFillingPlaces[i] = make(map[int]bool)
	}
	for i := 0; i < 12; i++ {
		status.FuelBalance += config.Deliveries[i]
		server.FinishBuildingNewObjects(&status, config, i)
		server.CountMaintenancePrice(&status, config)
		server.ChooseStrategy(&status, config, i)
		for localTime := 0; localTime < month; localTime++ {
			server.ManageTankers(&status, config, i)
			server.ProcessStations(&status, config)
		}
		server.Pay(&status, &logger)
		fmt.Println(status.MoneyBalance, status.FuelBalance)
	}

	router.HandleFunc("/config", server.GetConfig).Methods("POST")

	//router.HandleFunc("/api/search", search).Methods("GET", "OPTIONS")
	//router.HandleFunc("/api/add-to-cart", addToCart).Methods("POST", "OPTIONS")
	//router.HandleFunc("/api/cart", getCart).Methods("GET", "OPTIONS")
	//router.HandleFunc("/api/delete-from-cart", deleteFromCart).Methods("POST", "OPTIONS")
	spa := spaHandler{staticPath: "public", indexPath: "index.html"}

	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}
