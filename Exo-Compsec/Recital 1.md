
**A DDoS attack prevents you from connecting to your bank website. Which of the security
properties will this impact?**

Availability.

**What is the difference between authenticity and integrity**

Authenticity: le message est bien envoyé par la bonne personne.
Integrity: le message est bien celui envoyé à l'origine.

**What are the basic elements of a threat model?**

- **assets** : ce qu'on veut protéger
- **vulnerabilities** : les faiblesses dans le système qui peuvent être exploitées
- **attacks** : une action qui exploite une vulnérabilité pour exécuter une menace
- **controls**: atténuer ou supprimer une vulnérabilité.

**Image that you have important data on your laptop. You place a tracking chip inside it in a
tamper resistant enclosure. Is this a cost effective way to protect your laptop? Relate your
answer to the different types of defences one could employ to protect their assets. Make
sure to include your assumptions.**

Non. On se protège mal contre la menace, qui est le vol de l'ordinateur et l'extraction de données. La question est moins de savoir où l'ordinateur se trouve que de savoir que les données sont illisibles par quiconque, et qu'elles peuvent être retrouvées en cas de perte/vol.

Au lieu de ça, on pourrait mieux protéger l'accès à l'ordinateur (bureau plus sécurisé), chiffrer le disque et faire des backups régulières.

**ARP allows address translation between IP and MAC addresses.**

- **How many MAC addresses are allocated to each manufacturer for their use (assuming
one prefix each).** 

La moitié (les 3 premiers octets déterminent l'organisation).

- **Which of the two address spaces would be exhausted first, MAC or IPv4? Give the
(approximate) difference.** 

IPv4, $2^32$ vs $2^48$ pour MAC.

**Recall encapsulation. Imagine that a packet of 10 bytes needs to go through 3 layers of
the stack before it is transmitted to another machine. Each layer added 10bytes of header
and 2 bytes of footer.**

- **What is the size of the packet that is transmitted?**

10 bytes (data) + 3 x 12 bytes = 46 bytes

- **Imagine that the original 10byte packet is fragmented into two. Now what is the total size
(in bytes) of transmitting that original packet?**

10 bytes (data) + 2 x 3 x 12 bytes = 82 bytes

**NAT is useful to ease the exhaustion pressure on the IPv4 address space. It can also hide
the internal information of a private network from external observers. Give at least one type
of information that could be prevented from being observed? Give reasons why this is good
to protect.**

One would be which computer is making the request (we don't know which student is asking for a specific file inside the school for instance). 

**Imagine you want to divert internet traffic to your own knock-off bank website that is a
duplicate of the original bank website.**

- **What could you do to divert the traffic from the real website to yours (assume that
certificates or other forms of authentication are not present)?**

Advertise false IP ranges using BGP.
Inside a private network, answer an ARP request with my mac address for any IP.

- **Would this be a stealthy attack, or would it be traceable?**

Traceable, we know who sent the message.

**Imagine that an IDS has been trained to detect website-X (that serves malware) and the IDS has a TPR = 95.99% and an FPR = 15%. Suppose that website-X is very popular and 50% of all website visits that the IDS observes are to it. What is the probability that when the IDS detects a visit to website-X it is correct? Show your intermediate steps.**

true positive rate : 95.99 P(websiteX detected / websiteX)
false positive rate : 15 P(websiteX detected / !websiteX)

$95.99*0.5 + 15*0.5 = 0.55$ 
50%

P (websiteX / websiteX detected) = P(websiteX and websiteX detected) / P(website X detected)

P(websiteX and websiteX detected) = P(websiteX detected / websiteX) * P(websiteX) = $0.5 * 95.99 = 86$