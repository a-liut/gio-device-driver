package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gio-device-driver/cmd/pkg/api"
	"gio-device-driver/cmd/pkg/model"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	checkVariables()

	host := flag.String("host", "localhost", "IP address of the current host")
	port := flag.Int("port", 8080, "port to be used")

	flag.Parse()

	go registerService(*host, *port)

	log.Printf("Server started on port %d", *port)

	router := api.NewRouter()

	p := fmt.Sprintf(":%d", *port)

	log.Fatal(http.ListenAndServe(p, router))
}

func registerService(host string, port int) {
	callbackUuid, err := registerCallback(host, port)
	if err != nil {
		panic(err)
	}

	log.Printf("Callback UUID: %s", callbackUuid)
}

func registerCallback(host string, port int) (string, error) {
	fogNodeHost := os.Getenv("FOG_NODE_HOST")
	fogNodePort := os.Getenv("FOG_NODE_PORT")

	u := fmt.Sprintf("http://%s:%s", fogNodeHost, fogNodePort)
	log.Printf("FogNode URL: %s\n", u)

	_, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	callbackUrl := fmt.Sprintf("http://%s:%d%s", host, port, api.CallbackEndpointPath)
	callbackData := struct {
		Url string `json:"url"`
	}{
		Url: callbackUrl,
	}

	dataJson, _ := json.Marshal(callbackData)

	registrationUrl := fmt.Sprintf("%s/callbacks", u)
	registrationResp, err := http.Post(registrationUrl, "application/json", bytes.NewBuffer(dataJson))
	if err != nil {
		return "", err
	}

	var message model.ApiResponse
	err = json.NewDecoder(registrationResp.Body).Decode(&message)
	if err != nil {
		return "", err
	}

	// Return the UUID
	return message.Message, nil
}

func checkVariables() {
	if fogNodeHost := os.Getenv("FOG_NODE_HOST"); fogNodeHost == "" {
		panic("FOG_NODE_HOST not set.")
	}
	if fogNodePort := os.Getenv("FOG_NODE_PORT"); fogNodePort == "" {
		panic("FOG_NODE_PORT not set.")
	}

	if DeviceServiceHost := os.Getenv("DEVICE_SERVICE_HOST"); DeviceServiceHost == "" {
		panic("DEVICE_SERVICE_HOST not set.")
	}
	if DeviceServicePort := os.Getenv("DEVICE_SERVICE_PORT"); DeviceServicePort == "" {
		panic("DEVICE_SERVICE_PORT not set.")
	}
}
