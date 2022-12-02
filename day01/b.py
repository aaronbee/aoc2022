#!/usr/bin/env python3

import sys


f = open(sys.argv[1])
v = [[]]
cur = v[0]
for line in f:
    line = line.strip()
    if line == "":
        cur = []
        v.append(cur)
        continue
    cur.append(int(line))

summed = [sum(series) for series in v]
print(sum(list(reversed(sorted(summed)))[0:3]))