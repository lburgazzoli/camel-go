	kcat \
		-t "${CAMEL_KAFKA_TOPIC}" \
		-b "${CAMEL_KAFKA_BROKER}" \
#		-X security.protocol=SASL_SSL
#		-X sasl.mechanisms=PLAIN \
#		-X sasl.username="${CAMEL_KAFKA_USER}" \
#		-X sasl.password="${CAMEL_KAFKA_PASSWORD}" -C