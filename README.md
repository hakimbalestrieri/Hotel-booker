# PRR 2021

## Introduction  / Labo 1

**Auteurs :** Hakim Balestrieri & Alexandre Mottier

**Dates :** 30 septembre 2021 au 31 octobre 2021

L’hôtel dispose de plusieurs chambres qui peuvent être réservées une ou plusieurs nuits consécutives par des clients. L’application fonctionne en mode client-serveur. Un client s’identifie par son nom sur le serveur de l’hôtel puis établi une réservation en précisant le jour, le numéro de chambre et le nombre de nuits souhaité de sa réservation. Il peut également obtenir une liste de l’occupation des chambres de l’hôtel un jour donné. Il peut aussi demander un numéro de chambre disponible à partir d’un jour donné pour un nombre de nuits fixés. Il n'est pas demandé de gérer la persistance des réservations au-delà de l'exécution du serveur. 



## Introduction  / Labo 2

**Auteurs :** Hakim Balestrieri & Alexandre Mottier

**Dates :** 21 octobre 2021 au 28 novembre 2021 

Partager  des  variables  parmi  un  ensemble  de  processus  est  un  problème  qui  peut  se  résoudre  par  le biais  d'un  algorithme  d'exclusion mutuelle.  Dans ce  laboratoire,  nous  allons  utiliser l’algorithme de Lamport comme algorithme d’exclusion mutuelle pour gérer l’accès à la section critique. Il  s’agit  de  reprendre  la  petite  application  de  réservation  développée  dans  le  1er  laboratoire  et  de  la faire  fonctionner  avec  plusieurs  serveurs,  sans  serveur  central.  Chaque  client  détient  une  copie  des variables partagées et doit obtenir l’accès à des sections critiques pour y apporter des modifications et diffuser celles-ci sur les autres clients. 


## Introduction  / Labo 3

**Auteurs :** Hakim Balestrieri

**Dates :** 28 novembre au 23 décembre 2021 

TODO


## Initialisation

Pour cloner notre projet en SSH : 

```
git clone git@github.com:alex-mottier/PRR-Labo3-Balestrieri.git
```

Pour build notre serveur tcp : 

```
go.exe build tcpServer.go
```

Pour build notre client tcp : 

```
go.exe build tcpClient.go
```



## Connexion

Pour instancier notre serveur tcp, nous avons mis en place dans notre fichier de config 3 serveurs à disposition (les ports sont configurés dans le fichier de configuration aussi) : 

```
.\tcpServer 1
```

```
.\tcpServer 2
```

```
.\tcpServer 3
```

Pour instancier un ou plusieurs clients tcp, il faut ajouter en argument le numéro du serveur sur lequel se connecter : 

```
.\tcpClient 1
```



## Business

L'hôtel dispose de 10 chambres qui ont chacune 31 dates disponibles, les clients connectés au serveur peuvent utiliser les commandes suivantes :

| Commande  | Arguments                          | Description                                                  | Remarque                                   |
| --------- | ---------------------------------- | ------------------------------------------------------------ | ------------------------------------------ |
| /username | USERNAME                           | Ajoute un nom d'utilisateur au client                        | Est obligatoire pour /reserve ou /rooms    |
| /reserve  | NUMERO_CHAMBRE, DATE, NOMBRE_NUITS | Permet de reserver une chambre                               | /username est obligatoire à exécuter avant |
| /rooms    | DATE, (optionel) NOMBRE_NUITS      | Permet de lister les chambres, si l'argument optionel est ajouté, le serveur enverra en réponse une chambre disponible | /username est obligatoire à exécuter avant |
| /quit     |                                    | Déconnecte le client du serveur                              |                                            |



## Network

Un protocole a été créé pour gérer les messages de Lamport afin de gérer les accès à la section critique : 

| Commande | Arguments     | Description |
| -------- | ------------- | ----------- |
| /req     | Lamport Clock |             |
| /ack     | Lamport Clock |             |
| /rel     | Lamport Clock |             |

Un autre protocole a été créé pour gérer les messages d'update afin de gérer les updates des clients et des rooms sur chaque serveur : 

| Commande    | Arguments                                    | Description                                                  |
| ----------- | -------------------------------------------- | ------------------------------------------------------------ |
| /upt_client | USERNAME                                     | Permet de mettre à jour le dictionnaire de clients de chaque hôtel afin d'éviter les doublons |
| /upt_rooms  | NUMERO_CHAMBRE, DATE, NOMBRE_NUITS, USERNAME | Permet de mettre à jour la map de chambres disponible dans l'hôtel et attribuer la chambre au bon client |

La commande **/ready** permet de savoir que tous les autres serveurs sont prêts.

La commande **/hello** permet de savoir si le serveur qu'on contacte nous écoute bien.

## Scénario de tests

Pas eu le temps de faire, comme Lamport n'est pas terminer.



## Remarques

#### Labo1

Les corrections du labo1 données par Raphaël Racine ont été prises en compte et les problèmes ont été corrigés durant ce laboratoire.

Les tests du labo1 ont été corrigés et améliorés pour correspondre à un scénario d'utilisation.

#### Labo2

Il peut y avoir des problèmes de connexion entre les serveurs, il suffit de relancer les serveurs pour que ça fonctionne. 

L'algorithme de Lamport n'est pas 100% fonctionnel. Parfois cela fonctionne bien, mais la plupart du temps non. Cela est surement dû a un problème entre les channels.

La partie REQ et ACK de l'algorithme fonctionne, mais parfois le serveur ne parvient pas à envoyer le message de release aux autres serveurs.
