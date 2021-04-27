package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/SubrotoRoy/event-consumer/model"
)

//cache is a map with string key and model.PriceDerails as value.
var cache = map[string]model.PriceDetail{}

//GetFuelPrice returns the per litre cost of fuel for a given city
func GetFuelPrice(city string) float64 {
	wg := &sync.WaitGroup{}
	m := &sync.RWMutex{}
	cacheCh := make(chan float64)
	apiCh := make(chan float64)

	resultCh := make(chan float64)

	wg.Add(2)

	//fetches data from cache
	go func(city string, wg *sync.WaitGroup, m *sync.RWMutex, cacheCh chan float64) {
		if p, ok := getFromCache(city, m); ok {
			cacheCh <- p
		}
		wg.Done()
	}(city, wg, m, cacheCh)

	//fetches data from API
	go func(city string, wg *sync.WaitGroup, m *sync.RWMutex, apiCh chan float64) {
		if p, ok := getFromAPI(city, m); ok {
			apiCh <- p
		} else {
			log.Println("Price Could not be fetched for the given city")
			apiCh <- p
		}
		wg.Done()
	}(city, wg, m, apiCh)

	//co-ordinates the responses from cache and API,
	//whichever one is received first that is retured as result
	go func(cachech, apiCh chan float64) {
		select {
		case b := <-cacheCh:
			log.Println("From cache")
			log.Println(b)
			resultCh <- b
			<-apiCh
		case b := <-apiCh:
			log.Println("From api")
			log.Println(b)
			resultCh <- b
		}

	}(cacheCh, apiCh)

	return <-resultCh
}

//getFromCache queries the cache for the fuel cost
func getFromCache(city string, m *sync.RWMutex) (float64, bool) {

	//acquiring read lock to query the cache to avoid any race condition
	m.RLock()
	price, ok := cache[city]
	m.RUnlock()

	if time.Since(price.CreatedDate).Hours() >= 24 {
		return 0.00, false
	}

	return price.Price, ok
}

//getFromAPI makes a call to the external API for the fuel cost
func getFromAPI(city string, m *sync.RWMutex) (float64, bool) {
	priceList := model.PriceResponse{}
	getPriceURL := os.Getenv("PRICEURL")
	response, err := doGetAPICall(getPriceURL + city)

	//if there is a failure to get response from API
	if err != nil {
		log.Println("Error encountered while doing get call", err)
		return 0.0, false
	}

	//If there is no failure then Unmarshalling the response received into struct
	json.Unmarshal(response, &priceList)
	if len(priceList.Prices) == 0 {
		return 0.0, false
	}

	//As per API contract the first structure is the latest price
	petrolPrice := priceList.Prices[0].PetrolPrice

	//converting string valued price to float64
	price, _ := strconv.ParseFloat(petrolPrice, 64)
	priceDetail := model.PriceDetail{}
	priceDetail.Price = price

	//Converting the PriceDate as string to time.Time
	const layout = "2006-01-02 15:04:00.0"
	priceDetail.CreatedDate, err = time.Parse(layout, priceList.Prices[0].PriceDate)
	if err != nil {
		log.Println("Error encountered while getting time of price. Using current time")
		priceDetail.CreatedDate = time.Now()
	}

	//acquiring write lock to update the cache to avoid race condition
	m.Lock()
	cache[city] = priceDetail
	m.Unlock()

	return price, true
}

//doGetAPICall performs simple GET call to the provided URL
func doGetAPICall(url string) ([]byte, error) {
	var contents []byte
	resp, err := http.Get(url)
	if err != nil {
		return contents, err
	}
	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}
	return []byte{}, errors.New("API call returned status" + strconv.Itoa(resp.StatusCode))
}
