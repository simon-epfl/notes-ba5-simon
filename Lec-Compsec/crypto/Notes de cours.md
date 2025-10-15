## Symmetric encryption

**Perfect secrecy:** $|P(E(k, m_1) = c) - P(E(k, m_2) = c)| <= epsilon$

**Stream ciphers** : OTP but use a pseudorandom key rather than a really random key. The key will be generated from a key seed using a Pseudo-Random Generator (PRG).
Problèmes :
1. ⚠️ **Two-time pad attack**
    Si on réutilise la même clé (seed) pour deux messages → même flot pseudo-aléatoire → vulnérable comme le _two-time pad_.
    Donc la clé ne doit **jamais être réutilisée**.
2. ⚠️ **Malleabilité**
    Comme dans le OTP, un attaquant peut modifier des bits du ciphertext pour modifier de façon prévisible les bits du message déchiffré.

**Block ciphers** : A block cipher with parameters 𝑘 and ℓ is a pair of deterministic algorithms (𝐸, 𝐷) such that:
Encryption 𝐸 ∶ ${0, 1}^k times {0, 1}^l arrow {0, 1}^l$
Decryption 𝐷 ∶ ${0, 1}^k times {0, 1}^l arrow {0, 1}^l$

DES:  ℓ = 64 bits, 𝑘 = 56 bits
3DES: ℓ = 64, 𝑘 = 168

AES: ℓ = 128, 𝑘 = 128, 192, 256
• Goal: replace 3DES which is too slow (3DES is 3 times as slow as DES)

Le texte clair est organisé dans une **matrice de 4x4 octets** (appelée _state_).

| b0 | b4 | b8  | b12 |
| b1 | b5 | b9  | b13 |
| b2 | b6 | b10 | b14 |
| b3 | b7 | b11 | b15 |

Chaque “round” (il y en a 10 pour AES-128) applique ces opérations :
1. **SubBytes** 🧬
    → Chaque octet est remplacé par un autre selon une table fixe (_S-box_)
    → Ajoute de la **non-linéarité** : impossible de prédire simplement la sortie.
2. **ShiftRows** 🔄
    → Les lignes de la matrice sont décalées vers la gauche d’un certain nombre de positions.
    → Mélange les octets horizontalement.
3. **MixColumns** 🧪
    → Chaque colonne est multipliée par une matrice fixe (dans le corps de Galois).
    → Mélange les octets verticalement (diffusion).
    ⚠️ Cette étape n’est **pas appliquée au dernier tour**.
4. **AddRoundKey** ➕
    → On “mélange” (XOR) la matrice avec une **sous-clé** dérivée de la clé principale.
    → Cette clé change à chaque tour grâce à la **key expansion**.

#### Résumé

• Frequency analysis as a cryptanalysis attack on classic encryption
• Importance of randomness in cryptography
• Stream ciphers
• simple and efficient symmetric encryption schemes
• use a random IV to thwart two-time pad attacks
• subject to malleability attacks
• Block ciphers - use AES not DES
• CBC mode is more secure that ECB but less resilient to packets loss
• CTR mode more secure than ECB and parallelisable
• Keep up to date with cryptanalytic advances and standards.
Modern symmetric encryption also guarantees authenticity → no malleability
• Do not implement crypto lightly - use public reference implementations

## Asymmetric Encryption

TODO: 
- regarder video sur SHA-256 intuitvement
- vidéo about HMAC/PMAC +  https://www.youtube.com/watch?v=gZiBYDX9Fpo
- comprendre pourquoi RSA elle-même est pas secure "The “textbook RSA signature” scheme does not provide existential unforgeability"
##### **Merkle–Damgård (MD)**
1. **But : construire un hash pour un message de taille arbitraire**
    - La construction transforme un message de n’importe quelle longueur en un **hash fixe**.
2. **Découpage en blocs et padding**
    - Le message M est découpé en blocs m_1, m_2, \dots, m_t.
    - On ajoute un **padding (PB)** pour compléter le dernier bloc, souvent incluant la longueur du message.
3. **Fonction de compression**
    - Une **fonction de compression** h prend un **état interne** et un **bloc** et produit un nouvel état interne h : $\mathcal{T} \times \mathcal{X} \to \mathcal{T}$       
    - L’état initial est l’**IV** (Initial Value).    
4. **Chaînage des blocs**
    - On applique h bloc par bloc :
        $h_1 = h(IV, m_1), \quad h_2 = h(h_1, m_2), \dots, h_t = h(h_{t-1}, m_t)$
    - La dernière sortie h_t devient le **hash final** H(M).
