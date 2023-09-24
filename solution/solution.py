#!/usr/bin/python3

import argparse
import random


def swap(data: list[int], a: int, b: int):
    data[a], data[b] = data[b], data[a]
    print(data)


def _merge_inner(data: list[int], start: int, end: int):
    """
    _merge_inner is the actual implementation of merge sort
    """
    if start == end or start == end - 1:
        return
    cut = int((start + end) / 2)
    _merge_inner(data, start, cut)
    _merge_inner(data, cut, end)
    i = start

    # this makes copies
    left = data[start:cut]
    right = data[cut:end]

    for i in range(start, end):
        didSwap = False
        if len(left) > 0 and (len(right) == 0 or left[0] < right[0]):
            if data[i] != left[0]:
                data[i] = left[0]
                didSwap = True
            left = left[1:]
        else:
            if data[i] != right[0]:
                data[i] = right[0]
                didSwap = True
            right = right[1:]
        i += 1
        if didSwap:
            print(data)


def merge(data: list[int]):
    """
    merge implements merge sort. It is a convenience method over _merge_inner

    This one is a little weird because merge sort requires temp space. That
    means that the data slice will contain missing or double entries.
    """
    _merge_inner(data, 0, len(data))


def bubble(data: list[int]):
    """
    bubble implements bubble sort

    This implementation takes advantage of the fact that each cycle sets the
    final value of the last element of the considered range.
    """
    n = len(data)
    for s in range(1, n):
        for j in range(0, n - s):
            if data[j] > data[j + 1]:
                swap(data, j, j + 1)


def insertion(data: list[int]):
    """
    insertion runs insertion sort, as inspired by Romanian folk dancers
    """
    for i in range(1, len(data)):
        for j in range(i - 1, -1, -1):
            if data[i] < data[j]:
                swap(data, i, j)
                i = j
            else:
                break


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog="mike_sort",
        description="sorts and emits stuff",
    )
    parser.add_argument("-c", "--count", default=64, type=int)
    parser.add_argument(
        "algorithm",
        default="bubble",
        choices=["bubble", "insertion", "merge"],
    )
    args = parser.parse_args()
    random.seed()
    data = list(range(args.count, 0, -1))
    random.shuffle(data)
    print(data)
    if args.algorithm == "bubble":
        bubble(data)
    elif args.algorithm == "insertion":
        insertion(data)
    elif args.algorithm == "merge":
        merge(data)
    else:
        raise RuntimeError("unknown algorithm: " + args.algorithm)
