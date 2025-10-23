**Write down the code for an RM macro ‘if R1 > R2 then goto Ij ’. The macro must leave all registers unchanged after its execution. Assume a predefined goto macro.**

```
1.

2. DECJZ(R1, z)
3. DECJZ(R2, j)
4. GOTO 2
```

**Give a simple recursive definition of a sequence coding function N∗ → N, based on the pairing function in the slides.**

$<> = 0$, $<x, <y>> = 2^x 3^(<y>)$   

**Suppose we allow register machine to have an unbounded number of registers, but each register is finite (e.g. 32 bits) – like current computer memory. With no changes to the instruction set, are these machines still Turing powerful? Why or why not?**

They would not be, because we can only specify a finite number of registers.

**Suppose now that we add a form of indirect addressing. For example, we might say that the register operand of an instruction can now be either i, as before, meaning Ri, or (i), meaning RRi . Does that help?**

No, because we are still limited in the number of accessible registers because $i$ is bounded.

**Why aren’t Turing machines affected by this issue? Can you adapt ideas from TMs to solve this issue for RM?**

because TMs can access any cell. we would need another way of specifying which register to use (using an unlimited variable for instance).

**The barber shaves all and only the men who do not shave themselves. Does the barber shave himself?**
$arrow$ yes. but then he is not **only** shaving the men who do not shave themselves!
$arrow$ no. but then he is not shaving **all** the men who do not shave themselves!

**The set of sets that are not members of themselves.**

Problème pour savoir si l'ensemble lui-même doit être compris dedans.

**The smallest natural number not definable in under eleven words.**


