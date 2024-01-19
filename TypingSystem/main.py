from GoType import *
import math

# Code
for i in range(2, 100):
    prime: bool = True

    for j in range(2, i):
        if i%j == 0:
            prime = False
    
    if prime:
        print(i)

print(2)
for i in range(3, 100, 2):
    prime: bool = True
    
    for j in range(3, i, 2):
        if i%j == 0:
            prime = False
    
    if prime:
        print(i)

print(2)
for i in range(3, 100, 2):
    prime: bool = True
    
    for j in range(3, math.sqrt(i)+1, 2):
        if i%j == 0:
            prime = False
    
    if prime:
        print(i)