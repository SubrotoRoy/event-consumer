# event-consumer
event-consumer service reads events from a kafka topic. Once it gets a true and a false event then it processes them to determine the time interval between them and determines
then amount to be paid.

## Environment Variables required
  Variable Name: BROKER
  
  Variable Description: Required to get the broker url of kafka
  
  
  Variable Name: TOPIC
  
  Variable Description: Required to get the kafka topic name
  
  
  Variable Name: PRICEURL

  Variable Description: set it as https://mfapps.indiatimes.com/ET_Calculators/oilprice.htm?citystate= 
  
  This is a freely available API from where one can fetch the fuel price per litre.

  
## Installation
1. Clone the repository in your GOPATH

2. Install kafka, start the zookeeper, and start the broker.

3. Execute the below commands :
    NOTE: replace the parameters present in <> with actual values.
    
    i.  set BROKER=<broker_url>
    
    ii. set TOPIC=<topic_name>
    
    iii.set PRICEURL=https://mfapps.indiatimes.com/ET_Calculators/oilprice.htm?citystate=
        
    
4. Navigate to the repository's directory in your local and run the command

    go run main.go
 
 The application would start running and be ready to accept messages form the kafka topic
