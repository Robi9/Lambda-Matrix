#!/bin/bash

URL=$1

# Invocação Go
curl -H "Content-Type: application/json" \
-X POST \
-d '{ "id" : "1",  "ordem" : 10, "linguagem" : "go", "horario" : "", "tempo" : ""}' \
${URL}/multgo

# Invocação Python
curl -H "Content-Type: application/json" \
-X POST \
-d '{ "id" : "100",  "ordem" : 10, "linguagem" : "python", "horario" : "", "tempo" : ""}' \
${URL}/multpython
