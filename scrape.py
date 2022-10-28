#!/usr/bin/python3
import json
import requests
import sys

f = open("import.sql", "w+")

page = 0
users = 0
while True:
    req = requests.get("https://mee6.xyz/api/plugins/levels/leaderboard/{}".format(sys.argv[1]), params={"page": page}).json()
    page += 1
    for user in req["players"]:
        users += 1
        id = user["id"]
        xp = user["xp"]
        level = user["level"]
        f.write(f"INSERT INTO levels (id, xp) VALUES ({id}, {xp});\n")
        print(f"\r Current user level: {level} ({users} total users) ", end='')
        sleep(1)

f.close()
