import random


def swap(data, a, b):
    data[a], data[b] = data[b], data[a]
    print(data)


fart = list(range(100))
swap(fart, 3, 6)
random.shuffle(fart)
print(fart)
swapcount = -1
while swapcount != 0:
    swapcount = 0
    for ix in range(len(fart) - 1):
        if fart[ix] > fart[ix + 1]:
            swap(fart, ix, ix + 1)
            swapcount += 1
