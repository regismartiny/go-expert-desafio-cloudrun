package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/regismartiny/desafio-cloudrun/configs"
	"github.com/regismartiny/desafio-cloudrun/internal/viacep"
	"github.com/regismartiny/desafio-cloudrun/internal/weatherapi"
)

type HandlerData struct {
	ViacepClient     *viacep.Client
	WeatherapiClient *weatherapi.Client
}

func main() {

	config, _ := configs.LoadConfig(".")

	h := &HandlerData{
		ViacepClient:     getViaCepClient(config.ViaCepAPIBaseURL, config.ViaCepAPIToken),
		WeatherapiClient: getWeatherClient(config.WeatherAPIBaseURL, config.WeatherAPIToken),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /temperatura/{cep}", http.HandlerFunc(h.handleGet))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Erro:", err)
	}
}

func getViaCepClient(baseURLStr string, apiToken string) *viacep.Client {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		log.Fatal(err)
	}
	return viacep.NewClient(baseURL, apiToken)
}

func getWeatherClient(baseURLStr string, apiToken string) *weatherapi.Client {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		log.Fatal(err)
	}
	return weatherapi.NewClient(baseURL, apiToken)
}

func (h *HandlerData) handleGet(w http.ResponseWriter, r *http.Request) {

	cep := r.PathValue("cep")

	ctx := context.Background()

	//Viacep

	adrressInfo, err := h.ViacepClient.GetAddressInfo(&ctx, cep)
	if err != nil {
		log.Println(err)
	}

	log.Println("Address info")
	log.Println(adrressInfo)

	cidade := adrressInfo.Localidade

	log.Println("Cidade: " + cidade)

	//WeatherAPI

	weatherInfo, err := h.WeatherapiClient.GetWeatherInfo(&ctx, cidade)
	if err != nil {
		log.Println(err)
	}

	log.Println("Weather info")
	log.Println(weatherInfo)

	json := fmt.Sprintf("{ \"temp_C\": %f, \"temp_F\": %f, \"temp_K\": %f }",
		weatherInfo.Current.TempC,
		weatherInfo.Current.TempF,
		convertCelsiusToKelvin(weatherInfo.Current.TempC))

	log.Println(json)
	w.Write([]byte(json))
}

func convertCelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}
