
Recall that in the lectures we showed that the class of regular languages was closed
under union, sequential composition, and Kleene closure.

**The complement of a language $L$, written $L$, is every string not in L, i.e. $Sigma^* backslash L$. Show that the regular languages are closed under complement.**

Est-ce que le complément d'un langage est aussi un langage régulier ? oui, on inverse les states finaux / non finaux.

**Using (a), or other arguments, show that the regular languages are closed under
intersection. The intersection of two languages L1 ∩ L2 is the set of all strings that
are in both L1 and L2, that is ${w | w ∈ L_1 ∧ w ∈ L_2}$**

$$ L_1 sect L_2 = L_1 union L_2 - (L_1 sect Sigma^* backslash L_2) - L_2 sect Sigma^* backslash L_1 $$
TODO

**An all-NFA is a variant of an NFA where a string w is only accepted if all states
reached on word w are final, i.e. $δ^*(q_0, w) ⊆ F$. Show that there is an all-NFA
that recognises a language if and only if it is regular.** 

C'est-à-dire qu'au lieu d'avoir un test "en parallèle" et de regarder si au moins une des solutions conduit à un state final, on le fait sur tous. 

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

Myhill-Nerode:
- si on prend $0^i 1^j 2^(m - n - 1), i >= j > 1$, et qu'on ajoute $z = 2$ à la fin, on obtient toujours un mot valide.
-  si on prend $0^i 1^j 2^(m - n - 2), i >= j > 1$, et qu'on ajoute $z = 2$ à la fin, on obtient toujours un mot invalide.

