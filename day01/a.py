#!/usr/bin/env python3

import sys
import itertools


def main():
    f = open(sys.argv[1])
    it = (line.strip() for line in f)

    emptyLineCount = 0

    def emptyLineCounter(k: str):
        nonlocal emptyLineCount
        if k == "":
            emptyLineCount += 1
        return emptyLineCount

    sums = [
        sum(int(l) for l in g if l != "")
        for _, g in itertools.groupby(it, emptyLineCounter)
    ]

    print(max(sums))


if __name__ == "__main__":
    main()
