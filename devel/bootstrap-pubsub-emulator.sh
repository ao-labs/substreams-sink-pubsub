#!/usr/bin/env bash

set -euo pipefail

emulator_addr="http://${PUBSUB_EMULATOR_HOST:-"localhost:8085"}"
project="${PROJECT:-"acme"}"
topic="${TOPIC:-"dev-topic"}"
subscription="${SUBSCRIPTION:-"dev-topic-sub"}"

output=`curl -sS -X PUT "${emulator_addr}/v1/projects/${project}/topics/${topic}"`
if [[ $output == *"\"code\""* ]]; then
  if [[ $output != *"ALREADY_EXISTS"* ]]; then
    echo $output
    exit 1
  fi
fi

curl \
    -X PUT \
    -H "Content-Type: application/json" \
    -d "{\"topic\": \"projects/${project}/topics/${topic}\"}" \
    "${emulator_addr}/v1/projects/${project}/subscriptions/${subscription}"
