### Regular languages

Recall that in the lectures we showed that the class of regular languages was closed
under union, sequential composition, and Kleene closure.

**The complement of a language $L$, written $L$, is every string not in L, i.e. $Sigma^* backslash L$. Show that the regular languages are closed under complement.**

Est-ce que le complément d'un langage est aussi un langage régulier ? oui, on inverse les states finaux / non finaux.

⚠️ Correction: ça ne marche **que** si l'automate est déterministe. Donc on doit bien préciser que chaque NFA peut être converti en DFA (en gardant le même langage), donc l'argument de flipper les states marche quand même.

**Using (a), or other arguments, show that the regular languages are closed under
intersection. The intersection of two languages L1 ∩ L2 is the set of all strings that
are in both L1 and L2, that is ${w | w ∈ L_1 ∧ w ∈ L_2}$**

$$ L_1 \cap L_2 = \overline{\overline{L_1} \cup \overline{L_2}} $$

⚠️ Correction : penser à de Morgan... on pouvait bien exprimer l'intersection en fonction de l'union.

**An all-NFA is a variant of an NFA where a string w is only accepted if all states
reached on word w are final, i.e. $δ^*(q_0, w) ⊆ F$. Show that there is an all-NFA
that recognises a language if and only if it is regular.** 

C'est-à-dire qu'au lieu d'avoir un test "en parallèle" et de regarder si au moins une des solutions conduit à un state final, on le fait sur tous.

⚠️ Correction: on a un NFA qui a un langage $L_1$. Pour que $w in L_1, exists "path"$ avec au moins un final state.
Si on prend le complémentaire, $\overline{L_1}$, pour que $w \in \overline{L_1}, \forall \text{path}$, le state n'est pas final via la machine NFA. On peut donc maintenant construire une all-NFA qui reconnaît $\overline{L_1}$ en inversant les states finaux $arrow.l.r$ rejects. On a donc un langage régulier !

- Every regular language has an NFA → flip to an all‑NFA → complement stays regular.
- Every all‑NFA’s language is regular because it’s the complement of an NFA’s language.

**Prove that $L = {0^n 1^m 2^(m−n) | m ≥ n ≥ 0}$ is not regular. You may use any
of the three methods used in lectures: the Pumping Lemma, the Myhill-Nerode
theorem, or using (or proving) closure properties of the regular languages to reduce
the problem to a known non-regular language. As an extension exercise, you may
wish to try multiple methods.**

Pumping lemma:
- $0^(p-1) 1^p$
- donc soit $x = 0^(p-1)$, $y = 1$, $z=1^(p-1) 2$
- soit $x = 0^i$, $y = 0^(p - i - 1) 1, z=1^(p-1) 2$
$arrow$ dans tous les cas, on ne peut pas pump $y$, car on aura toujours le 2 à la fin qui dépend du nombre de $0$ et de $1$. et que le 2 est forcément en dehors de $x y$.

⚠️ ok, bien tester tous les découpages.

Myhill-Nerode:
- si on prend $0^i 1^j 2^(m - n - 1), i >= j > 1$, et qu'on ajoute $z = 2$ à la fin, on obtient toujours un mot valide.
-  si on prend $0^i 1^j 2^(m - n - 2), i >= j > 1$, et qu'on ajoute $z = 2$ à la fin, on obtient toujours un mot invalide.
### Context-free languages

Consider the CFG G: $S → a S | a S b S | ε$

**Informally characterise L(G).**

- a
- aaba
- aba
- aaaaaaaaababa

une suite de a et b qui commence forcément par un a, et pas de b collés entre eux.

**Show that it is ambiguous by finding a string for which you can construct two
parse trees and two leftmost derivations.**

aaba

S $arrow$ $a S b S$ $arrow$ $a (a S) b S$ $arrow$ $a a b a S$ $arrow$ $a a b a$
S $arrow$ $a S$ $arrow$ $a (a S b S)$ $arrow$ $a a b (a S)$ $arrow$ $a a b a$ 

**Find an unambiguous grammar for L(G).**

$T arrow a S b | a$
$S arrow T S | epsilon$

**Give a pushdown automaton that recognises L(G).**

![[image-19.png]]

**Some $w ∈ {a, b}^∗$ have unique parse trees in G. What are those strings w? Give an efficient test to tell whether w has this property. The test “try all parse trees to see how many yield w” is not adequately efficient.**

Est-ce qu'il y a plus d'un a à la suite avant un b.

⚠️ Correction:
```python
balance = 0
ambiguous = False
for i, sym in enumerate(w):
	if sym == 'a':
		balance += 1
	else: # b
		if balance == 0:
			reject() # invalid
		balance -= 1
		# A b can attach to any of the open a's → ambiguous
		if balance >= 1:
			ambiguous = True
```

### Reductions for undecidability

In lectures, we considered the Halting Problem H, the Looping Problem L, and the Uniform Halting Problem UH . Now we consider the Universal Looping problem UL: given a machine M , does M loop on all inputs R?

**Show, by reduction from L, that UL is undecidable.**

On construit une machine $M prime$ qui ignore son input, écrit $w$ sur le tape et appelle $M$. Si $M$ loop sur $w$, alors $M prime$ loop sur tous les inputs aussi.

**Show, by constructing a suitable machine, that UL is co-semi-decidable. (Hint: interleaving.)**

