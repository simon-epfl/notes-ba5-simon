### Finite Automata

Language : toutes les chaînes que l'automate A valide $\cal{L} (A)$
Alphabet $\Sigma$

![[Pasted image 20250923161419.png]]

les états finaux sont décidés à l'avance et si on "tombe" à la fin de la chaîne dans un autre état on considère que la chaîne est rejetée.

#### Finite Automata Déterministe (DFA)

La fonction de transition donne exactement un state de destination pour chaque symbole d'entrée (déterministe)

mais on peut aussi avoir... 
$arrow$ fonction de transition **partielle** : certains symboles ne pointent vers rien
$arrow$ fonction de transition **non-déterministe** : certains symboles pointent vers deux states

#### Finite Automata Non-Déterministe (NFA)

La différence entre un automate **déterministe** et **non-déterministe** est que la fonction de transition est $delta : Q times Sigma arrow cal(P) (Q)$ où $cal(P) (Q)$ est le power-set de l'ensemble des states $Q$.

L'automata non déterministe va tester en parallèle toutes les possibilités jusqu'à trouver une suite d'état qui donne un état final. Dans ce cours, on parle d'**angelic** non-determinism, donc on regarde si **au moins une** des suites d'état arrive à un état final.

#### DFA = NFA

Chaque automate non déterministe peut être transformé en automate déterministe (et vice-versa mais c'est beaucoup plus simple). 

Au lieu de fonctionner avec l'ensemble des states $Q$, on utilise le power-set des states $Q$, $cal(P) (Q)$. Au début on part de l'ensemble des states de départ, puis à chaque fois on rajoute l'ensemble des states sur lesquels on pourrait être d'une façon ou d'une autre (c'est le **subset construction**).

![[Pasted image 20250923173253.png|384x271]]

$epsilon$ est la chaîne vide. on appelle les transitions $q_x arrow_epsilon q_y$ des silent moves.

$E(q)$ ($epsilon$-closure) sont l'ensemble des states atteignables par depuis $q$ par des silent moves.
on a donc comme fonction de transition :
$$ delta_D (S, a) = E ( union_(s in S) " " delta_A (s, a)) $$
### Propriétés des langages

- les langages **réguliers** sont ceux qui peuvent être exprimés avec des DFA, NFA et $epsilon$-NFA
- l'**union** de deux langages $L_1$ et $L_2$, écrit $L_1 union L_2$ est le langage qui inclus toutes les chaînes de $L_1$ et de $L_2$.
- on dit que l'union de deux langages réguliers est **fermée**, c'est-à-dire que l'union est elle-même régulière et peut-être exprimée avec un DFA, NFA ou $epsilon$-NFA.

> [!tip] **Preuve du dernier point**:
> - on part de $N_1 = (Q_1, sigma, delta_1, s_1, F_1)$ et $N_2 = (Q_2, sigma, delta_2, s_2, F_2)$
> - on crée un nouveau state de départ $s_0$.
> - on crée une $epsilon$-transition entre $s_0 arrow s_1$ et $s_0 arrow s_2$.
> - on considère les states finaux $F = F_1 union F_2$

- la **composition séquentielle** de deux langages, écrite $L_1 L_2$ est le langage des chaînes qui consistent en une chaîne de $L_1$ suivie d'une chaîne de $L_2$. Elle est aussi **fermée** (on créé une $epsilon$-transition entre chaque state final de $L_1$ vers le state de départ de $L_2$).

- on appelle **Kleene-closure** d'un langage $L$, écrite $L^*$, le langage des chaînes qui consistent de zéro ou plus chaînes de $L$. Elle est aussi fermée, on ajoute une $epsilon$-transition du state final vers le state de départ.
### Regex

$R^*$ : zéro ou plus chaînes de $R$
$R^+ = R compose R^*$ : une ou plus chaînes de $R$
$R?$ : zéro ou une fois une chaîne de $R$

Chaque langage régulier peut être représenté sous la forme d'une expression régulière.

### G-NFA (Generalised Non-deterministic Finite Automata)

- les GNFA sont des NFA dont les transitions sont des expression régulières et non des symboles.
- elles n'ont qu'un state final.
- on veut que tous les states soient reliés entre eux sauf l'initial state qui n'a pas de transition entrante et le state final qui n'a pas de transition sortante.

**DFA $arrow$ GNFA**

- on ajoute un nouveau state initial qu'on connecte à l'ancien state initial via une $epsilon$-transition
- on ajoute un nouveau state final auquel on connecte tous les anciens states finaux via des $epsilon$-transitions
- si deux states $q_0$ et $q_1$ ont deux transitions entre eux $q_0 arrow^a q_1$ et $q_0 arrow^b q_1$, on les remplace par $q_0 arrow^(a union b) q_1$.
- comme on veut que tous les states soient reliés entre eux, on ajoute des $emptyset$-transitions. 

**DFA $arrow$ RE**

On veut éliminer tous les states intérieurs de la GNFA un par un. Quand ils sont tous partis, seuls les states initial/final vont rester avec une unique transition entre eux. Son label sera notre expression régulière.