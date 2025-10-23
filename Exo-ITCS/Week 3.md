
**Consider the language $L = {w w | w ∈ {a, b}^∗}$. Show that it is not regular.**

Intuitivement, le langage a besoin de garder en mémoire l'ancien mot pour vérifier qu'il y a bien le même ensuite.

On peut utiliser le pumping lemma. On dit qu'il y a un mot $w$ de taille $p$, et on montre que si on split le mot en $x, y, z$ alors on ne va pas pouvoir répéter la partie du milieu autant de fois qu'on veut, puis qu'on ajoutera des caractères répétés à la fin du premier mot, qui n'apparaîtront pas à la fin du second mot.

**Is the language $L′ = {w w | w ∈ {a, b}^∗ ∧ |w| < 4}$ regular? Why or why not?**

le langage est régulier. en effet, le fait que la taille du mot est borné fait qu'on peut garder trace avec un NFA facilement des anciens caractères.
pour le prouver, on peut construire le NFA ?

**We saw in lectures that the language L is not context-free. Show that, by contrast, its complement is context-free.**

We can define the following grammar :

$S -> A | A B | B A | B$
$A -> a | b A b | a A b | a A a | b A a$
$B -> b | b B b | a B b | a B a | b B a$

----

The following questions/comments are intended as prompts for discussion. Of course, you can ask/discuss about anything. Some of these topics we’ve touched on in discussion in lectures – this is an opportunity to think about them a bit more.

**Suppose we augmented a finite automaton with an additional feature: a single mutable variable. On each transition, it may read from or write to this variable.**

- **If the variable is of type Σ (that is, it can contain one alphabet symbol), what class of languages can such automata recognise?**

⚠️ If it can only hold a **finite** symbol (like one letter from the alphabet), the machine is still finite → **regular languages**.

- **If the variable is of type N, what class of languages can they recognise?**

they can recognize some context-free languages, but not all.

⚠️ If it can hold an **integer** and perform arithmetic, it becomes very powerful — potentially **Turing-complete**, since all reasonable data structures can be encoded as one number.

- **A queue automaton is like a pushdown automaton, but with a queue (fifo) instead of a stack (lifo). These can recognise any Turing-recognisable language. Now suppose we instead defined a pushdown automaton with an additional stack. What class of languages can these 2-stack PDAs recognise?**

they can recognize some turing-recognisable languages, but not all.

⚠️ A **queue automaton** (FIFO) is equivalent in power to a Turing machine.
You can **simulate a queue using two stacks**:
- One stack acts as temporary storage to reverse order when needed.
Therefore, a **2-stack PDA** can recognize any **Turing-recognisable language**.