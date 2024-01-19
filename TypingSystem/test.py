from GoType import *
# Code
for i in range(2, 100):
    prime: bool = True
    for j in range(2, 100):
        if i%j == 0:
            prime = False
    if prime:
        print(i)
print(i)