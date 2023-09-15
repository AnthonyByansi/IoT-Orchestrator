package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Define a struct for IoT devices.
type Device struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

var devices []Device

func main() {
	// Create a new router.
	router := mux.NewRouter()

	// Define API routes.
	router.HandleFunc("/devices", GetDevices).Methods("GET")
	router.HandleFunc("/devices/{id}", GetDevice).Methods("GET")
	router.HandleFunc("/devices", CreateDevice).Methods("POST")
	router.HandleFunc("/devices/{id}", UpdateDevice).Methods("PUT")
	router.HandleFunc("/devices/{id}", DeleteDevice).Methods("DELETE")

	// Start the server.
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", router)
}

// GetDevices returns a list of all devices.
func GetDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// GetDevice returns a specific device by ID.
func GetDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range devices {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Device not found", http.StatusNotFound)
}

// CreateDevice adds a new device to the list.
func CreateDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var device Device
	_ = json.NewDecoder(r.Body).Decode(&device)
	device.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	device.LastUpdated = time.Now()
	devices = append(devices, device)
	json.NewEncoder(w).Encode(device)
}

// UpdateDevice updates an existing device by ID.
func UpdateDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range devices {
		if item.ID == params["id"] {
			var updatedDevice Device
			_ = json.NewDecoder(r.Body).Decode(&updatedDevice)
			updatedDevice.ID = item.ID
			updatedDevice.LastUpdated = time.Now()
			devices[index] = updatedDevice
			json.NewEncoder(w).Encode(updatedDevice)
			return
		}
	}
	http.Error(w, "Device not found", http.StatusNotFound)
}

// DeleteDevice removes a device by ID.
func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range devices {
		if item.ID == params["id"] {
			devices = append(devices[:index], devices[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Device not found", http.StatusNotFound)
}
