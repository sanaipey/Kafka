#!/bin/bash
BROKERS="localhost:9092"
docker run -it --rm \
    confluentinc/cp-kafkacat:4.1.3 \
    kafkacat -t myTopic -b "${BROKERS}" \
    -X "debug=consumer" \
    -C
