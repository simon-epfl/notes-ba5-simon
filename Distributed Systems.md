
## Styles architecturaux

### Layered Style

- **Couche Interface (Application)** : interaction avec l’utilisateur ou d’autres applications (ex. : navigateur web).
- **Couche Logique (Processing)** : contient la logique métier / traitements (ex. : moteur d’application, calculs).
- **Couche Données (Data)** : stockage et gestion des données (ex. : base de données).

**Avantages** :
- Modulaire → facile à maintenir, modifier une couche sans casser les autres.
- Séparation claire des responsabilités.

### Service-Oriented Style (SOA & Microservices)

- **Object-based** : objets distribués qui communiquent via appels de procédures (RPC).
- **Microservices** : application découpée en petites unités autonomes, chacune ayant une responsabilité précise.
- **Resource-based** : chaque ressource est accessible via une API uniforme (ex. REST).

**Avantages** :

- Grande flexibilité → chaque service peut évoluer indépendamment.
- Facile à déployer et mettre à l’échelle (scalabilité).

**Exemple** : Netflix ou Amazon : chaque fonctionnalité (paiement, recherche, recommandation) est un microservice.

### Publish-Subscribe Style (Pub/Sub)

- **Fonctionnement** :    
    - Les producteurs publient des événements ou messages.
    - Les consommateurs s’abonnent à certains types d’événements.
    - Un **broker** (ou un espace de données partagé) distribue les messages aux abonnés.
- **Deux variantes** :
    - **Event-based** : notifications envoyées dès qu’un événement survient.
    - **Shared data-space** : dépôt commun où producteurs déposent des données et où consommateurs viennent les lire.
- **Avantages** :
    - Forte **découplage** : les producteurs ne connaissent pas les consommateurs.
    - Très adapté aux systèmes distribués et scalables.
- **Exemple** :
    - **Kafka, RabbitMQ** (systèmes de messagerie).
    - Réseaux sociaux (un utilisateur publie un post, ses abonnés le reçoivent automatiquement).
