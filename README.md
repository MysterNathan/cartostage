# CartoStage

CartoStage est une application web de cartographie interactive des offres de stage, permettant aux étudiants, enseignants et tuteurs entreprise de gérer et visualiser les stages disponibles.

> Projet full-stack développé de A à Z — du provisioning de l'infrastructure à l'application, en passant par la base de données et la CI/CD.

![Carte interactive](./docs/screenshots/map.png)

---

## Stack technique

### Frontend
![Next.js](https://img.shields.io/badge/Next.js-000000?style=flat&logo=nextdotjs)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=flat&logo=typescript&logoColor=white)
![Leaflet](https://img.shields.io/badge/Leaflet-199900?style=flat&logo=leaflet&logoColor=white)

### Backend
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white)

### Infrastructure & DevOps
![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=flat&logo=kubernetes&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=flat&logo=githubactions&logoColor=white)
![Traefik](https://img.shields.io/badge/Traefik-24A1C1?style=flat&logo=traefikproxy&logoColor=white)

---

## Architecture

L'application repose sur une architecture **microservices**, chaque service étant indépendant, déployé dans son propre conteneur et exposé via **Traefik** (routing, load balancing, middlewares d'authentification).

![Schema Architecture](./docs/screenshots/archimc.png)


| Service | Rôle |
|---|---|
| `auth-service` | Authentification, émission et validation des tokens JWT |
| `stage-service` | Gestion des offres de stage et des données cartographiques |
| `enterprise-service` | Profil et gestion des tuteurs entreprise |
| `student-service` | Profil et dashboard étudiant |
| `teacher-service` | Profil et dashboard enseignant |

> Les trois services utilisateurs (`enterprise`, `student`, `teacher`) partagent une structure similaire. Ce découpage est volontaire : il évite tout couplage inter-domaine et permet de faire évoluer chaque profil indépendamment.

---

## Fonctionnalités

### Carte interactive
- Visualisation des offres de stage géolocalisées via **OpenStreetMap / Leaflet**
- Filtrage par filière, secteur, disponibilité
- Fiche détaillée par stage au clic

![Carte interactive](./docs/screenshots/map.png)

### Espaces utilisateurs
Chaque profil dispose d'un espace dédié avec un dashboard adapté à son contexte :

| Profil | Fonctionnalités |
|---|---|
| **Étudiant** | Consultation des stages, suivi de candidature |
| **Enseignant** | Suivi des étudiants, gestion des fiches de stage |
| **Tuteur entreprise** | Gestion des offres, suivi des stagiaires |

![Dashboard](./docs/screenshots/dashboard.png)

### Formulaires
- Rédaction et soumission de fiches de stage
- Validation et gestion du cycle de vie des documents

---

## CI/CD

La pipeline est déclenchée sur `push` vers `main`/`develop` et sur `pull_request` vers `main`, avec une **détection intelligente des changements** pour ne rebuilder que ce qui est impacté.

![Schema github workflow](./docs/screenshots/workflowgh.png)

**Points notables :**
- Si le répertoire `shared/` backend est modifié, les tests sont déclenchés sur **l'ensemble des services**
- Le service `auth-service` est défini comme **blocking** : ses tests en échec bloquent tout déploiement
- Les images sont poussées vers un **registry Docker privé** hébergé dans le cluster
- Le déploiement se fait via `kubectl rollout` directement depuis le **runner self-hosted** sur l'infrastructure

---

## Tests

La couverture de tests est une démarche active tout au long du développement, avec une couverture élevée sur l'ensemble des services.

| Service | Type | Couverture |
|---|---|---|
| `auth-service` | Unitaires | Élevée |
| `stage-service` | Unitaires | Élevée (~90%+) |
| `enterprise-service` | Unitaires | Élevée |
| `student-service` | Unitaires | Élevée |
| `teacher-service` | Unitaires | Élevée |

Les tests sont exécutés dans un environnement isolé avec une base **PostgreSQL de test** provisionnée directement dans la CI via Docker.

> Les tests d'intégration sont prévus dans les prochaines évolutions du projet.

---

## Infrastructure

Ce projet est hébergé sur une infrastructure **Kubernetes on-premise** provisionnée via **Ansible**, documentée dans un repository dédié :

**[MysterNathan/k8s-infra](https://github.com/MysterNathan/k8s-infra)** — Cluster Kubernetes bare-metal avec stack d'observabilité (Prometheus, Grafana), sécurité runtime (Falco), stockage distribué (Longhorn) et opérateur PostgreSQL HA.

---

## Perspectives

- Complétion des dashboards par profil (étudiant, enseignant, tuteur)
- Gestion complète des stages côté frontend — CRUD et amélioration des formulaires
- Portail entreprise pour la soumission autonome d'offres de stage
- Extension de la couverture de tests en parallèle du développement
- Mise en place de tests d'intégration
- Enrichissement de la CI : lint, rapport de couverture, checks qualité

---

## Screenshots

| | |
|---|---|
| ![Login](./docs/screenshots/login.png) | ![Carte](./docs/screenshots/map.png) |
| ![Dashboard](./docs/screenshots/dashboard.png) | ![Formulaire](./docs/screenshots/form.png) |
