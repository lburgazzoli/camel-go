mosquitto_pub -h test.mosquitto.org -p 1883 -t camel-iot -m '{ "source": "sensor-1", "data": "foo" }'
mosquitto_pub -h test.mosquitto.org -p 1883 -t camel-iot -m '{ "source": "sensor-1", "data": "foo" }' | jq .

kubectl -n redpanda-system exec -ti redpanda-0 -c redpanda -- \
  rpk topic consume --brokers redpanda-0.redpanda.redpanda-system.svc.cluster.local.:9093 -f json near far
