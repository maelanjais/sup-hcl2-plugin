# Guide de Démonstration : Orchestration Multi-VM avec HCL2

Ce guide détaille un scénario complet d'utilisation du plugin HCL2 pour Sup, mettant en oeuvre un cluster de 4 machines virtuelles gérées avec Tart. L'objectif est de démontrer la capacité du plugin à gérer des configurations complexes et l'exécution parallèle.

## Prérequis matériels et logiciels

Avant de commencer, assurez-vous de disposer des éléments suivants :
- Un Mac avec Apple Silicon (pour Tart).
- Le binaire `sup` (version modifiée) compilé.
- Le binaire `sup-hcl2-parser` compilé et disponible dans le dossier courant ou dans le PATH.
- L'outil Tart installé (`brew install cirruslabs/cli/tart`).

## Étape 1 : Initialisation de l'infrastructure locale

Nous allons créer un cluster de 4 serveurs Ubuntu identiques.

```bash
# Récupération de l'image officielle Ubuntu
tart pull ghcr.io/cirruslabs/ubuntu:latest

# Clonage des 4 noeuds du cluster
for i in {1..4}; do tart clone ghcr.io/cirruslabs/ubuntu:latest "node-$i"; done

# Lancement des machines en mode arrière-plan
for i in {1..4}; do nohup tart run "node-$i" > /dev/null 2>&1 & done
```

Attendez environ 30 secondes que les systèmes démarrent, puis récupérez les adresses IP :

```bash
echo "--- Inventaire du cluster ---"
for i in {1..4}; do echo "node-$i : $(tart ip "node-$i")"; done
```

## Étape 2 : Configuration des accès SSH

Sup utilise les clés SSH pour communiquer. Nous devons injecter notre clé publique dans les nouvelles VMs (le mot de passe par défaut de l'image est `admin`).

```bash
# Remplacez les adresses IP par celles obtenues à l'étape précédente
for ip in 192.168.64.10 192.168.64.11 192.168.64.12 192.168.64.13; do
  ssh-copy-id -o StrictHostKeyChecking=no admin@$ip
done
```

## Étape 3 : Définition de la configuration HCL2

Créez un fichier nommé `Supfile-complex.hcl`. Ce fichier utilise des variables d'environnement, plusieurs réseaux et des commandes structurées.

```hcl
version = "0.5"

# Variables globales pour le déploiement
env {
  APP_NAME = "demo-app"
  LOG_PATH = "/var/log/demo.log"
}

# Définition du cluster de noeuds
network "cluster" {
  hosts = [
    "admin@192.168.64.10",
    "admin@192.168.64.11",
    "admin@192.168.64.12",
    "admin@192.168.64.13"
  ]
}

# Commande d'inspection système
command "check-system" {
  desc = "Vérifie les ressources système sur tous les noeuds"
  run  = "uname -a && free -h && df -h /"
}

# Commande de simulation de déploiement
command "prepare-node" {
  desc = "Prépare les répertoires et fichiers de log"
  run  = "sudo mkdir -p /opt/${APP_NAME} && sudo touch ${LOG_PATH} && ls -l ${LOG_PATH}"
}

# Commande locale (exécutée sur le Mac)
command "local-audit" {
  desc  = "Génère un rapport local de début de déploiement"
  local = "echo 'Début de l'audit à :' $(date) > deployment.audit"
}

# Target regroupant plusieurs étapes
target "full-setup" {
  commands = ["local-audit", "check-system", "prepare-node"]
}
```

## Étape 4 : Exécution de la démonstration

Lancez maintenant l'orchestration complète. Sup va charger le plugin HCL2, convertir la configuration et l'exécuter en parallèle sur les 4 machines.

```bash
./sup -f Supfile-complex.hcl --parser ./sup-hcl2-parser cluster full-setup
```

### Analyse du résultat attendu

Vous devriez voir dans votre terminal :
1. L'exécution de la commande locale (création du fichier `deployment.audit`).
2. Le flux de sortie des 4 machines s'affichant simultanément, préfixé par leur adresse IP.
3. La création des dossiers sur chaque VM avec les bons droits.

## Étape 5 : Nettoyage de l'environnement

Une fois la démonstration terminée, vous pouvez supprimer les machines virtuelles :

```bash
killall tart
for i in {1..4}; do tart delete "node-$i"; done
```

---

Ce scénario démontre que la transition vers HCL2 permet de conserver toute la puissance de Sup tout en bénéficiant d'un format de configuration plus moderne et structuré.
