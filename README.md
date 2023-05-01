# neobrain-todolist

neobrain-todolist est un projet de test technique qui consiste à réaliser une API en Golang permettant de gérer une "TODO list". L'API est conçue avec les endpoints suivants :

- Créer une tâche
- Supprimer une tâche
- Marquer une tâche comme complétée
- Récupérer la liste de toutes les tâches en filtrant les tâches complétées ou non
- Import et export des tâches sous un format CSV ou Excel

Le projet utilise Docker, docker-compose, MySQL pour la base de données et Golang pour l'API.

## Prérequis

Vous devez avoir installé Docker et docker-compose sur votre machine pour exécuter ce projet.

## Structure du projet

- `/cmd/api`: Contient le fichier `main.go` pour lancer l'API
- `/docker`: Contient les Dockerfiles pour l'API et la base de données
- `/internal`: Contient les fichiers sources de l'API en Golang, organisés en différents packages (api, config, database, utils)
- `/scripts/sql`: Contient le script d'initialisation de la base de données et de la table (`init.sql`)
- `/example`: Fichiers d'exemple pour l'import de tâches
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

   Cette commande va construire les images Docker nécessaires (Golang pour l'API et MySQL pour la base de données) et
   lancer les conteneurs. Le script d'initialisation de la base de données et de la table (`scripts/sql/init.sql`) sera
   automatiquement exécuté lors du lancement du conteneur MySQL.

## Utilisation de l'API

L'API sera accessible à l'adresse `http://localhost:8080`. Voici les endpoints disponibles :

### Créer une tâche

- Méthode : `POST`
- Endpoint : `/tasks`
- Description : Crée une nouvelle tâche en envoyant un JSON avec les propriétés `title` et `description` dans le corps
  de la requête.

### Supprimer une tâche

- Méthode : `DELETE`
- Endpoint : `/tasks/:id`
- Description : Supprime une tâche par son ID.

### Marquer une tâche comme complétée

- Méthode : `PUT`
- Endpoint : `/tasks/:id/complete`
- Description : Marque une tâche comme complétée par son ID.

### Récupérer la liste de toutes les tâches

- Méthode : `GET`
- Endpoint : `/tasks?completed=[true|false]`
- Description : Récupère la liste de toutes les tâches et permet de filtrer les tâches complétées (`true`) ou non
  complétées (`false`).

### Importer des tâches

- Méthode : `POST`
- Endpoint : `/tasks/import`
- Description : Importe des tâches à partir d'un fichier JSON. Le fichier doit être envoyé en tant que pièce jointe avec
  le paramètre `file`.

### Exporter des tâches

- Méthode : `POST`
- Endpoint : `/tasks/export`
- Description : Exporte les tâches vers un fichier JSON. Le fichier généré sera renvoyé en tant que réponse avec un
  téléchargement de fichier.

## Arrêt de l'application et nettoyage

Pour arrêter l'application et supprimer les conteneurs, les volumes et les réseaux, exécutez la commande suivante dans
le répertoire du projet :

```
docker-compose down
```