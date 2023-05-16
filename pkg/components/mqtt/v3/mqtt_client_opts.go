package v3

import mqtt "github.com/eclipse/paho.mqtt.golang"

type OptionFn func(*mqtt.ClientOptions)
