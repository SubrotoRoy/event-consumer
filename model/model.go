package model

import "time"

//Event struct for receiving messages from kafka
type Event struct {
	FuelLid   bool      `json:"fuellid"`
	City      string    `json:"city"`
	EntryTime time.Time `json:"-"`
}

//Price for accepting response of the external API
type Price struct {
	CityState   string `json:"cityState"`
	PetrolPrice string `json:"petrolPrice"`
	PriceDate   string `json:"priceDate"`
}

//PriceDetails struct is for creating channels and passing data around
type PriceDetail struct {
	Price       float64
	CreatedDate time.Time
}

//PriceResponse for accepting response of the external API
type PriceResponse struct {
	Prices []Price `json:"results"`
}
