package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const externalIPService string = "http://ipecho.net/plain"

var (
	mqttBroker = flag.String("mqttBroker", "", "MQTT broker URI (mandatory). E.g.:192.168.1.1:1883")
	topic      = flag.String("topic", "", "Topic where hub-ctrl messages will be received (mandatory)")
	user       = flag.String("user", "", "MQTT username")
	pwd        = flag.String("password", "", "MQTT password")
	period     = flag.Int("period", 3, "Periodic time to recheck the external IP address")
)

//Connect to the MQTT broker
func connectMQTT() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + *mqttBroker)

	if *user != "" && *pwd != "" {
		opts.SetUsername(*user).SetPassword(*pwd)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("%s", token.Error())
	}

	return client, nil
}

func init() {
	flag.Parse()
}

func main() {

	//Check command line parameters
	if *mqttBroker == "" || *topic == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	/*clientMQTT, err := connectMQTT()
	if err != nil {
		log.Fatalf("Error connecting to MQTT broker: %s", err)
	}

	log.Printf("Connected to MQTT broker at %s", *mqttBroker)*/

	for {

		resp, err := http.Get(externalIPService)
		if err != nil {
			log.Printf("Could not obtain the external IP address")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error retrieving HTTP body: %s", err)
			}
			ip := string(bodyBytes)
			log.Printf("External IP address is: %s", ip)
		}

		time.Sleep(10 * time.Second)

	}

}
