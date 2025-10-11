
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

## PgTeach

- `\q` quits the psql prompt and returns to the shell command line
- `\d` lists the tables stored in the database
- `\d+ table_name` gives details about the table
- `\i /absolute/path/to/filename` reads and executes SQL statements from the specified file
- `\ir relative/path/to/filename` same as the above, but the file path is relative to the one from which psql was invoked
- `\?` shows help on blackslash commands