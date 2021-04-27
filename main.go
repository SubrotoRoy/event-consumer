package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/SubrotoRoy/event-consumer/kafkaservice"
	"github.com/SubrotoRoy/event-consumer/model"
	"github.com/SubrotoRoy/event-consumer/service"
)

func main() {
	log.Println("Event-Consumer started")
	kafkaSvc := kafkaservice.NewKafkaService()

	eventSvc := service.NewEventManager()
	ctx := context.Background()

	var trueEvent, falseEvent *model.Event

	for {
		//reading message from kafka
		msg, err := kafkaSvc.ReadFromKafka(ctx)
		if err != nil {
			log.Println("Unable to read message from Kakfa. ERROR:", err)
			continue
		}
		log.Println("Received Message,", string(msg.Value), "at", msg.Time)

		//Unmarshalling kafka message to Event struct
		decodedMessage := model.Event{}
		err = json.Unmarshal(msg.Value, &decodedMessage)
		if err != nil {
			log.Println("Unable to unmarshall kafka message into struct. ERROR:", err)
			kafkaSvc.CommitMessage(ctx, msg)
			continue
		}

		//commiting read message
		kafkaSvc.CommitMessage(ctx, msg)

		decodedMessage.EntryTime = msg.Time

		//accepting true event only if there is no previously present true event and
		//accepting false event only if there is previously present true event
		if decodedMessage.FuelLid && trueEvent == nil {
			trueEvent = &decodedMessage
		} else if !decodedMessage.FuelLid && trueEvent != nil {
			falseEvent = &decodedMessage
		}

		//once both true event and false event is available then going for price calculation
		if trueEvent != nil && falseEvent != nil {
			result := eventSvc.GetPriceForEventSet(*trueEvent, *falseEvent)
			handleResult(result)
			trueEvent = nil
			falseEvent = nil
		}
	}
}

//handleResult takes care of the result generated. For now it is logging it to console.
func handleResult(result string) {
	log.Println("OUTPUT:", result)
}
