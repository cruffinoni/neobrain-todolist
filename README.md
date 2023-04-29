# neobrain-todolist

neobrain-todolist est un projet de test technique qui consiste à réaliser une API en Golang permettant de gérer une "TODO list". L'API est conçue avec les endpoints suivants :

- Créer une tâche
- Supprimer une tâche
- Marquer une tâche comme complétée
- Récupérer la liste de toutes les tâches en filtrant les tâches complétées ou non

Le projet utilise Docker, docker-compose, MySQL pour la base de données et Golang pour l'API.

## Prérequis

Vous devez avoir installé Docker et docker-compose sur votre machine pour exécuter ce projet.

## Structure du projet

- `/cmd/api`: Contient le fichier `main.go` pour lancer l'API
- `/docker`: Contient les Dockerfiles pour l'API et la base de données
- `/internal`: Contient les fichiers sources de l'API en Golang, organisés en différents packages (api, config, database, utils)
- `/scripts/sql`: Contient le script d'initialisation de la base de données et de la table (`init.sql`)
- `docker-compose.yml`: Fichier de configuration pour docker-compose
- `README.md`: Ce fichier, qui explique comment lancer le projet

## Mise en place de l'environnement

1. Clonez le dépôt Git sur votre machine :

   ```
   git clone https://github.com/user/neobrain-todolist.git
   ```

2. Accédez au répertoire du projet :

   ```
   cd neobrain-todolist
   ```

3. Lancez l'application en utilisant docker-compose :

   ```
   docker-compose up -d
   ```

   Cette commande va construire les images Docker nécessaires (Golang pour l'API et MySQL pour la base de données) et lancer les conteneurs. Le script d'initialisation de la base de données et de la table (`scripts/sql/init.sql`) sera automatiquement exécuté lors du lancement du conteneur MySQL.

## Utilisation de l'API

L'API sera accessible à l'adresse `http://localhost:8000`. Voici les endpoints disponibles :

- `POST /tasks` : Créer une nouvelle tâche (envoyez un JSON avec les propriétés `title` et `description` dans le corps de la requête)
- `DELETE /tasks/:id` : Supprimer une tâche par son ID
- `PUT /tasks/:id/complete` : Marquer une tâche comme complétée par son ID
- `GET /tasks?completed=[true|false]` : Récupérer la liste de toutes les tâches et filtrer les tâches complétées (`true`) ou non complétées (`false`)

## Arrêt de l'application et nettoyage

Pour arrêter l'application et supprimer les conteneurs, les volumes et les réseaux, exécutez la commande suivante dans le répertoire du projet :

```
docker-compose down
```