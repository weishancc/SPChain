#!/bin/bash

ORGNAME=$(sed 's/./\U&/2' <<< $(jq '.client.organization' connection-collector.json))
echo "${ORGNAME}"
