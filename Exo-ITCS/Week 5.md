**Let Even be the decision problem (N, {n : n is even}) and Odd be (N, {n : n is odd}). Give m-reductions between the two problems.**

$f(n) = n + 1$

**We remarked that (computable) predicates are (computable) functions. A function f : N → N is also a predicate: its characteristic predicate χf (x, y) is true iff y = f (x). Show that f is computable iff χf is computable.**

1. Pour calculer \(f(x)\), énumère les entiers \(y = 0, 1, 2, ...\)
2. Pour chacun, calcule \(\chi_f(x,y)\)
    - s’il vaut 1, c’est que \(y=f(x)\).
    - on s’arrête et on renvoie ce \(y\).

Ce procédé **termine toujours**, car \(f(x)\) est un **nombre bien défini** (la fonction est totale).

**Suppose that Q1 and Q2 are semidecidable queries over the same domain D. Show that Q1 ∪ Q2 and Q1 ∩ Q2 are semidecidable**

$Q_1 union Q_2$ --> si les deux ne s'arrêtent pas
$Q_1 sect Q_2$ --> si une des deux s'arrête pas

**Suppose that L is a decidable language over some alphabet Σ . Show that the language L∗ is decidable.**

systematically tries all possible cut points

**Following on from question 2: Given a computable predicate ψ(x, y), is it decidable whether ψ is the characteristic predicate of some function f ?**

no
the machine can never return