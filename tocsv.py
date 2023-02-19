#!/usr/bin/python3
import csv
import json
import sys

guild = sys.argv[1]

with open(f'./{guild}-levels.json') as jfile:
    j = json.load(jfile)

with open(f'{guild}-levels.csv', 'w', newline='') as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(['id', 'xp', 'guild'])
    for item in j:
        writer.writerow([item['id'], item['xp'], guild])