from GoType import *

for i in range(2, 1_000_000):
    prime: bool = True
    for j in range(2, i):
        if i%j == 0:
            prime = False
    if prime:
        print(i)