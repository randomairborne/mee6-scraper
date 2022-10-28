#!/usr/bin/python3
import json

ujson = open("results.json", "r")
levels = json.load(ujson)

f = open("import.sql", "w+")

for user in levels:
    id = user["id"]
    xp = user["xp"]
    f.write(f"INSERT INTO levels (id, xp) VALUES ({id}, {xp});\n")
f.close()

