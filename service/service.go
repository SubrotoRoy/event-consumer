package service

import (
	"strconv"

	"github.com/SubrotoRoy/event-consumer/model"
	"github.com/SubrotoRoy/event-consumer/utils"
)

type EventManager struct {
}

func NewEventManager() *EventManager {
	return &EventManager{}
}

//GetPriceForEventSet returns the total price for a given set of true and false events
func (e *EventManager) GetPriceForEventSet(trueEvent, falseEvent model.Event) string {
	if trueEvent.City != falseEvent.City {
		return "Invalid City names provided."
	}

	//fetching the per litle fuel price
	price := utils.GetFuelPrice(trueEvent.City)
	if price == 0.0 {
		return "Unable to fetch price"
	}

	//calculating the duration for which Fuel Lid was open
	dur := falseEvent.EntryTime.Sub(trueEvent.EntryTime)
	lidOpenTime := dur.Seconds()

	totalPrice := (lidOpenTime / 30) * price

	//returning total price upto 2 decimal places
	return strconv.FormatFloat(totalPrice, 'f', 2, 64)
}
