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
#### DFA $equiv$ NFA

Chaque automate non déterministe peut être transformé en automate déterministe (et vice-versa mais c'est beaucoup plus simple). 

Au lieu de fonctionner avec l'ensemble des states $Q$, on utilise le power-set des states $Q$, $cal(P) (Q)$. Au début on part de l'ensemble des states de départ, puis à chaque fois on rajoute l'ensemble des states sur lesquels on pourrait être d'une façon ou d'une autre (**subset construction**).

![[Pasted image 20250923173253.png|384x271]]

$epsilon$ est la chaîne vide. on appelle les transitions $q_x arrow_epsilon q_y$ des silent moves.

Soit $A = (Q_A, \Sigma, \delta_A, q_0, F_A)$ la machine non-déterministe **NFA**, et $q \in Q_A$ un état de $A$.
$E(q)$ (ou $epsilon$–closure de $q$) = ensemble des états atteignables depuis $q$ par des transitions $epsilon$.

On construit le **DFA** $D = (Q_D, \Sigma, \delta_D, q_0’, F_D)$.
- $Q_D = \mathcal{P}(Q_A)$
- $q_0’ = E(q_0)$
- états finaux : $F_D = \{ S \subseteq Q_A \mid S \cap F_A \neq \emptyset \}$
- fonction de transition $\boxed{\delta_D(S, a) = E\left( \bigcup_{s \in S} \delta_A(s, a) \right)}$

Idée : à partir de l’ensemble d’états S (un état du DFA), quand on lit le symbole a, on regarde  **toutes les transitions possibles** dans le NFA à partir de **chacun des états de** S, on regroupe ces destinations (union), puis on ajoute tous les états atteignables via des ε-transitions (ε-closure). C'est pour ça qu'on entoure tout avec un $E$ final.
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

- les GNFA sont des NFA dont les transitions sont des expressions régulières et non des symboles.
- elles n'ont qu'un state final.
- on veut que tous les states soient reliés entre eux sauf l'initial state qui n'a pas de transition entrante et le state final qui n'a pas de transition sortante.

**DFA $arrow$ GNFA**

![[image-4.png|498x278]]

- on ajoute un **nouveau state initial** qu'on connecte à l'ancien state initial via une $epsilon$-transition
- on ajoute un **nouveau unique state final** auquel on connecte tous les anciens states finaux via des $epsilon$-transitions
![[image-5.png|474x214]]

- si deux states $q_0$ et $q_1$ ont deux transitions entre eux $q_0 arrow^a q_1$ et $q_0 arrow^b q_1$, on les remplace par $q_0 arrow^(a union b) q_1$.
- comme on veut que tous les states soient reliés entre eux, on ajoute des $emptyset$-transitions. 

![[image-6.png|592x255]]

**DFA $arrow$ RE**

On veut éliminer tous les states intérieurs de la GNFA un par un. Quand ils sont tous partis, seuls les states initial/final vont rester avec une unique transition entre eux. Son label sera notre expression régulière.

![[image-7.png]]

### Langages non réguliers

Un automate fini n’a **qu’un nombre fini d’états** donc il **ne peut pas “compter” sans limite**.
Donc, tout langage qui **nécessite de compter ou de faire correspondre des éléments** (comme des parenthèses, des nombres égaux de symboles, etc.) **n’est pas régulier**.

Exemples :
- $L_1 = \{ a^n b^n \mid n \in \mathbb{N} \}$ (autant de $a$ que de $b$ et les $a$ viennent avant les $b$)
- $L_2 = \{ c^i a^j b^k \mid i = 1 \Rightarrow j = k + 1 \}$ (quand il y a **exactement un c**, on doit avoir **une a de plus que de b**)

#### Pumping lemma

si un langage est régulier, il doit respecter le pumping lemma : il existe un nombre $p$ (pumping length) tel que si $s$ est une chaîne de $A$ dont la taille est au moins $p$, alors $s$ peut être divisé en trois parties $x y z$ et satisfaire les conditions suivantes :
- $forall i >= 0, x y^i z in A$
- $|y| > 0$ (on ne pump pas la chaîne vide..)
- $|x y| <= p$ (notre pumping string se trouve dans les $p$ premiers symboles car après on repasse forcément par deux états différents comme on suppose qu'on a que $p$ états)

Soit $n = |Q|$, le **nombre d’états** de l’automate. On définit $p = n$.

L'idée c'est qu'on va pouvoir trouver une boucle dans la finite state machine ($y$), et la répéter plein de fois. Par contre $x$ et $z$ peuvent être vide.

⚠️ the converse of the pumping lemma **does not hold**: can’t be pumped $arrow.not.l.double$ not regular 
#### Left quotient d'un langage régulier

Définition : $w backslash L = {v | w v in L}$. 
Par exemple :$$c a backslash {c^i a^j b^k | i = 1 arrow.double j = k + 1} = {a^n b^n | n in NN }$$
(le set des suffixes qu'on peut ajouter à $c a$ pour qu'il reste dans le langage est l'ensemble des chaînes $a^n b^n, n in NN$)

Lemma : $L$ régulier implique $dot backslash L$ régulier.

#### Exact characterisation (Myhill-Nerode Theorem)

Soit un langage $L subset.eq Sigma^*$ et $x, y in Sigma^*$. S'il existe une chaîne suffixe $z$ telle que :
- $x z in L$
- $y z in.not L$
alors on dit que $x$ et $y$ sont **distinguishable** par $L$.

Si $x$ et $y$ ne sont pas distinguables par $L$, on dit que que $x equiv_L y$. C'est une **relation d'équivalence.**

**Myhill-Nerode Theorem** : un langage $L$ est régulier si et seulement si le nombre de classes d'équivalence $equiv_L$ est fini.

#### Context free languages

Une **context free grammar** (CFG) est un 4-tuple ($N, Sigma, P, S$) où :
- $N$ est un ensemble fini de variables (symboles non-terminaux)
- $Sigma$ est un ensemble fini de *terminals* (symboles terminaux)
- $P subset.eq N times (N union Sigma)^*$ est un ensemble fini de règles (ou *productions*). Par exemple $A arrow a B c$ ou $A arrow a | A a | b A b |$.
- $S in N$ est une variable de départ.

On dit qu'on fait une étape de dérivation $alpha A beta arrow.double_G alpha gamma beta$ (s'il existe une règle $A arrow gamma in P$).
Le langage d'une CFG est :
$$ cal(L) (G) = {w in Sigma^* | S arrow.double_G^* w}$$
Où $arrow.double_G^*$ est une closure réflexive (on autorise de ne faire aucune étape) et transitive (on peut enchaîner plusieurs étapes) de $arrow.double_G$. 

> [!tip] Un langage est **context-free** si et seulement il peut être reconnu par une context-free grammar.

**leftmost-derivation** : quand on remplace les caractères non terminaux de gauche à droite (on prend toujours celui le plus à gauche et on le remplace).

> [!tip] Les langages réguliers sont des langages context-free
> 
> $G = (V, \Sigma, R, S)$ qui génère exactement L.
> 
> Pour chaque état $q \in Q$, on introduit un non-terminal (une variable) $A_q$.
> Le symbole de départ est $S = A_{q_0}$
> Règles de production :
> - Pour chaque transition $\delta(q, a) = p$, ajouter : $A_q \to a A_p$
> - Si $q$ est un état final ($q \in F$), ajouter : $A_q \to \varepsilon$

Un **parse-tree** est un arbre qui monte comment générer une chaîne à partir d'un non-terminal.

![[image-8.png|184x277]]

Une grammaire est **ambigüe** si il existe au moins une chaîne $w$ de son langage qui peut être générée de deux façons différentes en termes d’arbre de dérivation (ou d’arbre syntaxique).

#### Pushdown automata

Une PDA est une $epsilon$-NFA avec un stack en plus.

![[image-9.png|591x248]]

Formellement c'est un 6-tuple: $Q, Σ, Γ, δ, q 0, F$.
Γ est le stack alphabet.
la fonction de transition $delta$ prend maintenant en param le stack symbol:
$$δ : Q times Σ ε × Γ ε →P(Q × Γ ε)$$
On se fiche du stack à la fin, la condition qu'on vérifie c'est toujours que le state final est un state accepté.

> [!info] Un langage est context-free si et seulement si il peut être reconnu par un pushdown automaton PDA.

#### CFG $arrow$ PDA

https://www.youtube.com/watch?v=GwS__G2M8mU

On veut avoir un PDA qui accepte les mêmes chaînes que celles qui sont générées par la context free grammar.

- On commence par créer un state $q_0$ et un state $q_1$, qu'on relie par une transition qui push $S$, la variable de départ, dans le stack.
- Ensuite on créé un state final (ici $q_2$), avec une transition avec un stack vide qui ne push rien.
- Ensuite, **pour chaque caractère terminal qu'on reçoit en entrée (la chaîne à vérifier)**, on ajoute dans notre $q_"loop"$ (ici $q_1$), une transition qui lit le terminal, l'enlève du stack, et ne push rien de plus.
- Ensuite on créé des boucles pour gérer chaque variable (voir schéma plus bas).

![[image-10.png]]- ![[image-11.png]]
#### PDA $arrow$ CFG

à vérifier, la PDA doit :
- avoir un unique state final
- vider son stack avant de finir
- n'avoir que des transitions qui push et pop un symbole (par les deux en même temps)
![[image-12.png]]

#### Pumping lemma pour les Context Free Languages

![[image-14.png]]

![[image-13.png]]

#### Chomsky Grammars

![[image-15.png]]

Généralisation des context-free grammars, l'élément à gauche peut contenir des terminaux mais **au moins** un non-terminal.

![[image-16.png]]

#### Emptiness

Est-ce qu'il existe un moyen de vérifier si :
- un **finite automaton** est vide ? --> oui avec DFS on parcourt les states et on regarde si on peut attendre un state final
- une **regex** est vide ? --> oui par induction
- un **context-free language** est vide ? oui on marque comme "generating" tous les terminals, et la chaîne vide. ensuite on regarde, pour chaque variable, si elle génère au moins un output avec un élément génératif. on répète jusqu'à ce qu'on ne rajoute plus de variable en "générating". On regarde si S est "generating".
#### Equivalence

Pour vérifier que deux DFAs sont équivalentes, on ne peut pas juste vérifier que $L_1 = L_2$ car les langages sont infinis. Par contre, on peut travailler sur les automates qui les reconnaissent.

L'idée c'est de construire la différence entre les deux ensembles :
$$L_3 = (L_1 sect L_2^C) union (L_1^C sect L_2) $$
On peut ensuite construire $L_3$ parce qu'on sait construire des machines à partir de l'union/intersection, puis vérifier si $L_3$ est vide. (et si c'est le cas, les deux sont égales).

Par contre, un algorithme ne peut pas vérifier si deux **CFGs** sont équivalentes (**problèmes indécidables**).

### Machines de Turing

On a donc vu les **DFAs/NFAs** (regular languages), les **CFG/PDAs** (context-free languages), mais certaines expressions comme ${a^n b^n c^n : n >= 0}$ n'est pas un CFL.

On a donc besoin d'un modèle égal à ce que peut vraiment faire un ordinateur.

Au lieu d'avoir un stack comme les PDAs, on a un ruban/tape.
- On écrit les inputs sur le ruban.
- On a une tape head, qui peut se déplacer de gauche à droite, lire et changer le contenu d'une cellule.
- Par défaut les cellules sont vides (il y a une infinité de cellules vides à droite).
- On peut aller à l'infini vers la droite.
- Si on essaye d'aller trop à gauche, on dit que la machine s'arrête, ou ne bouge pas, peu importe.

On a un input alphabet $Sigma$, un tape alphabet $Gamma$ (donc $Sigma subset.eq Gamma$ car on a l'input écrit sur le tape au début, et $union.sq in Gamma, union.sq in.not Sigma$, le caractère vide).

- On a un state $q_0 in Q$, le state initial.
- Un state $q_"accept"$ accepting state et un $q_"reject"$, rejecting state. Si la machine rentre dans un des deux, elle s'arrête (halts) et elle accepte/rejette l'input.

$$ delta : Q times Gamma ("elle prend un state et une cellule") \ arrow Q times Gamma times {"L", "R"} ("le state vers lequel on va, ce qu'on écrit dans la cellule et la prochaine direction") $$
$delta$ est une fonction **totale**, il n'y aura qu'une transition possible par state/cellule, et il y en a au moins une.
##### Terminologie

**Configuration** : une string dans $Gamma^* Q Gamma^*$. Les deux $Gamma^*$ sont les cellules à gauche et à droite de la cellule actuelle (et $q$ le state). donc l'étoile dit juste qu'on peut en avoir un ou plus. On considère que la cellule "en cours" est la première du $Gamma^*$ de droite.
du coup au début on a $q_0 i_1 i_2 ... i_n union.sq ..$ avec $i_j$ les inputs.

**Computation** : une séquence finite de configurations. $c_0$ est la config de démarrage, $c_(i+1)$ vient de $c_i$ après une transition. et comme la machine de turing est déterministe, on peut revenir en arrière. $c_(i+1)$ est "yielded from" $c_i$.

$$ L(M) = {w in Sigma^* : M "has an accepting computation on" w } $$
Un langage est dit **reconnaissable** si c'est le langage d'une machine de Turing. La machine doit s'arrêter au minimum sur les mots $w in L$.

Un langage est dit **décidable** s'il s'arrête sur tous les inputs.

##### Faire tourner une machine de Turing

On doit garder trace de:
- la chaîne en cours
- la position de la head tape
- le state

##### Encoder une machine

$arrow$ encoder une machine dans une chaîne de caractères. On suppose qu'on a un moyen de le faire, par exemple avec les fonctions de couplage (pairing functions). Une fonction de couplage est une fonction bijective qui permet d'encoder deux nombres naturels en un seul nombre naturel. La plus connue est la fonction de Cantor :
$$\langle x, y \rangle = \frac{(x + y)(x + y + 1)}{2} + y$$
Propriétés importantes :
- **Bijective** : chaque paire (x,y) correspond à un unique nombre naturel
- **Calculable** : on peut calculer ⟨x,y⟩ et extraire x,y à partir du résultat
- **Généralisation** : on peut encoder des n-uplets : ⟨x₁,x₂,x₃⟩ = ⟨⟨x₁,x₂⟩,x₃⟩

**Application à l'encodage de machines :**
1. **Instructions** :
    - INC(i) → ⟨0, i⟩
    - DECJZ(i,j) → ⟨1, i, j⟩ = ⟨1, ⟨i, j⟩⟩
2. **Programme** : ⟨I₀, I₁, ..., Iₙ₋₁⟩
3. **Registres** : ⟨R₀, R₁, ..., Rₘ₋₁⟩
4. **Machine complète** : ⟨Programme, Registres⟩

#### Halting problem

Supposons qu'on a un programme H qui détermine si un autre programme halts.

On créé une machine L qui, pour chaque input $I$ qu'il obtient, le donne à H et fait l'inverse de ce que H dit. Par exemple, si H dit que $I$ halts, alors $H$ va run forever, sinon halts.

La contradition arrive quand on donne le code de $D$ à $D$. Donc $D$ va donner son programme à $H$.
- disons que $H$ dit que $D$ va s'arrêter. mais en fait $D$ runs forever!
- disons que $H$ dit que $D$ runs forever. mais en fait $D$ s'arrête.

donc en fait même si on connaît l'input exact et le code de $D$, on ne peut pas créer un programme qui dit s'il s'arrête ou non.

**Turing transducer**

Supposons que tu as deux problèmes :

- **P₁** : difficile (par exemple, on ne sait pas s’il est décidable)
- **P₂** : un autre problème, qu’on ne sait pas encore classifier

Tu veux savoir si **P₂** est aussi difficile que **P₁**.  
Pour cela, tu construis une **machine** (ou une fonction) qui transforme _toute instance_ de **P₁** en une _instance équivalente_ de **P₂**. Résoudre P₂ te permettrait de résoudre P₁, **car tu peux transformer un cas de P₁ en un cas de P₂**.

C’est ça, une **réduction**.  
Et quand cette transformation est calculable par une machine de Turing (TM ou RM), on parle de **mapping reduction (ou many-one reduction)**.

##### Rice's theorem

![[image-18.png]]

https://www.youtube.com/watch?v=kr7n_3LpWhc&t=54s

Pour éviter de passer par une réduction, on peut passer par le Rice's theorem.

Par exemple : $"Three"_"TM" = {<M>, "M is a TM and " |L(M)| <= 3}$ 

- vérifier que c'est une propriété des TMs. si on a $L(M_1) = L(M_2)$, alors soit les deux appartiennent à Three ou aucune des deux. C'est donc une propriété des TMs. ✅
- est-ce que cette propriété est non triviale ? c-a-d est-ce qu'il y a une TM dans Three et une TM en dehors de Three ? Oui, une machine qui accepte tout, et une machine qui accepte qu'un mot. ✅

