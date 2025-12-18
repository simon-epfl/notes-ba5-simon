
## Basic SQL

`<>` : différent de

**Résultats ordonnés:** `ORDER BY col1, col2` $arrow$ si deux lignes ont la même valeur pour `col1`, alors elles sont triées en fonction de `col2`.

**Plus d'une table FROM** : les lignes des deux tables sont concaténés.
#### Jointures

Le `WHERE` permet de faire des jointures simples (`WHERE id = cust_id`).

On peut aussi les préciser explicitement :

```sql
SELECT Name, Balance
FROM Customer, Account
WHERE ID = CustID AND Balance < 0;

/* ou alors */
SELECT Name, Balance
FROM Customer JOIN Account ON ID = CustID
WHERE Balance < 0;
```

#### Renommer les attributs

```sql
SELECT C.Name AS CustName, A.Balance AS AccBal
FROM Customer C, Account A
WHERE A.CustID = C.CustID
AND A.Balance < 0;
```

#### Mettre à jour la db

```sql
ALTER TABLE name
	RENAME TO new_name ;
	RENAME column TO new_column ;
	ADD column type ;
	DROP column ;
	ALTER column
		TYPE type ;
		SET DEFAULT value ;
		DROP DEFAULT;
```

Supprimer les données d'une table : `TRUNCATE TABLE name;`
Enlever complètement une table (data + schéma) : `DROP TABLE name;`

`AS` optionnel
## Relational Algebra

Opérations :
- $pi$ projection (ne garder qu'une partie des colonnes) $arrow$ **SELECT**
- $sigma$ sélection (ne garder qu'une partie des lignes) $arrow$ **WHERE**
- $times$ produit (littéralement un produit cartésien des deux tables) $arrow$ **FROM**
- $rho$ remplacement (renommer un attribut)
- $union$ union
- $-$ différence

**Jointure** :$$ "Customer" join "Account" equiv pi_(X union Y) (sigma_("CustID" = "CustID'") ( "Customer" times rho_("CustID" arrow "CustID'") ("Account"))) $$ où $X$ est l'ensemble des attributs de Customer, et $Y$ l'ensemble des attributs de Account.

**Intersection** :
$$ R sect T = R - (R - S) $$

**Arity** : le nombre de colonnes dans la table.
**Cardinality** : le nombre de lignes.

**Division** :
$$ "Exams" div "DPT" = pi_"Student" ("Exams") - pi_"Student" (pi_"Student" ("Exams") times "DPT" - "Exams" ) $$

- On a une relation R(A, B) et une relation S(B).
- R ÷ S renvoie l’ensemble des valeurs de A telles que, pour toutes les valeurs de B présentes dans S, la paire (A, B) apparaît dans R.
- Autrement dit: on cherche les A “qui couvrent tous les B” de S.

## Propositional logic

$$alpha |= P equiv alpha (P) = t$$
> [!note] L’interprétation (ou valuation) $(\alpha)$ satisfait la formule $(\varphi)$.

$$ (\Sigma \vDash \varphi \quad \text{ssi pour toute interprétation } \alpha, \alpha \vDash \Sigma \implies \alpha \vDash \varphi) $$
![[image-21.png]]

On définit une interpréation avec une fonction de sémantique, qui nous permet de décrire ce qu'on veut dire par "Person" dans le contexte d'évaluation.

> [!note] $(\mathcal{I}, \nu \models \varphi)$ se lit :
> La formule $\varphi)$ est vraie (satisfaite) dans l’interprétation $(\mathcal{I})$ sous l’affectation $(\nu)$.

> [!note] Dans cette **interprétation** et avec cette **affectation**, la formule $(\varphi)$ est vraie:
> $$(\mathcal{I}, \nu) \text{ est un modèle de } \varphi \quad \text{si et seulement si} \quad \mathcal{I}, \nu \models \varphi)$$

![[image-22.png]]

> [!note] Une query est **safe**
> si elle n'utilise pas de constante (`WHERE salary > 4000`), si elle ne boucle pas indéfiniment, (**for any database that you could construct**). Déterminer si une query est safe est indécidable.

> [!note] Active domain
> 
> **adom**(_R_) = { all constants occuring in _R_ }
> 
> c'est l'ensembles des valeurs constantes se trouvant dans une table.
> 
> ![[image-23.png]]

> [!question] Safe calculus et Relational agebra?
> 
> Safe calculus:
> `∀x (Client(x) ∧ Pays(x, 'France')) ⇒ Nom(x, Ville(x))`
> 
> Relational algebra:
> `π_Nom, Ville (σ_Pays = 'France' (Clients))`

> [!tip] Relational Algebra $arrow$ Relational calculus
> 
> tout ce qu'on peut exprimer avec un, on peut l'exprimer avec l'autre.
> 
> P. exemple:
> $$E = \sigma_{\text{age} > 30}(\text{Personnes})$$
> $$\varphi(x) = \text{Personnes}(x) \land x.\text{age} > 30$$
> 
> **Relation de base**
> 
> Si $(E = R)$, une relation de base avec attributs $(x_1, \dots, x_n)$, $\varphi_R(x_1, \dots, x_n) = R(x_1, \dots, x_n)$
> 
> **Sélection**
> 
> Si $E = \sigma_{\theta}(E_1)$ (avec $\theta$ une condition) :
> $$\varphi_E(\mathbf{x}) = \varphi_{E_1}(\mathbf{x}) \land \theta(\mathbf{x})$$
> 
> > [!example]
> >   $\sigma_{\text{age} > 30}(Personnes)$ → $\varphi(x, a) = Personnes(x, a) \land a > 30$
> >   
> 
> **Projection**
> 
> Si $E = \pi_{A_1, \dots, A_k}(E_1)$ $arrow$ $\varphi_E(\mathbf{x}) = \exists \mathbf{y} \, \varphi_{E_1}(\mathbf{x}, \mathbf{y})$
> où les $\mathbf{y}$ sont les attributs qu’on a supprimés par la projection.
> > [!example]
> > $\pi_{\text{nom}}(Personnes)$ → $\varphi(n) = \exists a \, (Personnes(n, a))$
> 
> **Jointure naturelle (ou produit cartésien + sélection)**
> Si $E = E_1 \bowtie E_2$ $arrow$ $\varphi_E(\mathbf{x}, \mathbf{y}) = \varphi_{E_1}(\mathbf{x}) \land \varphi_{E_2}(\mathbf{y})$
> S’il y a des attributs en commun, on ajoute l’égalité sur ceux-ci.
> 
 

> [!tip] Le SQL utiise des **bags** : des sets avec duplication
> 
> On introduit $epsilon(R)$ qui enlève les duplicates de la relation $R$.
> L'intersection prend la multiplicté minimale.
> Si $a ∈_k R$ et $a ∈_n S$ alors $a in_(min k,n) in S sect R$.
> La différence soustrait la multiplicté jusqu'à potentiellement zéro.
> 
> ![[image-24.png]]
> 
> ![[image-25.png]]


![[image-26.png|517x388]]
![[image-27.png]]

**EXISTS** ( query ) is true if the result of query is **non-empty**

![[image-28.png]]
![[image-30.png]]
![[image-29.png]]![[image-31.png]]
## PgTeach

- `\q` quits the psql prompt and returns to the shell command line
- `\d` lists the tables stored in the database
- `\d+ table_name` gives details about the table
- `\i /absolute/path/to/filename` reads and executes SQL statements from the specified file
- `\ir relative/path/to/filename` same as the above, but the file path is relative to the one from which psql was invoked
- `\?` shows help on blackslash commands
