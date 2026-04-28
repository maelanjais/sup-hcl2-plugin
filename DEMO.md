# Guide de Démonstration : Orchestration Multi-VM avec HCL2

Ce guide détaille un scénario complet d'installation et d'utilisation du plugin HCL2 pour Sup. Nous allons mettre en œuvre un cluster de 4 machines virtuelles gérées avec Tart pour démontrer la puissance du plugin et de l'exécution parallèle.

## Étape 0 : Installation et Compilation

Avant de commencer, nous devons préparer les outils sur votre Mac.

### 1. Cloner les dépôts
```bash
mkdir -p ~/Documents/demo-sup && cd ~/Documents/demo-sup

# Cloner le moteur Sup (version modifiée pour supporter les plugins)
git clone https://github.com/maelanjais/sup.git

# Cloner le plugin HCL2
git clone https://github.com/maelanjais/sup-hcl2-plugin.git
```

### 2. Compiler les binaires
```bash
# Compilation de Sup
cd ~/Documents/demo-sup/sup
go build -o sup .

# Compilation du Plugin
cd ~/Documents/demo-sup/sup-hcl2-plugin
make build-plugin
```

### 3. Installer dans le PATH
Pour pouvoir utiliser les commandes de n'importe où, nous les installons dans `~/go/bin` :
```bash
mkdir -p ~/go/bin
cp ~/Documents/demo-sup/sup/sup ~/go/bin/
cp ~/Documents/demo-sup/sup-hcl2-plugin/sup-hcl2-parser ~/go/bin/
chmod +x ~/go/bin/sup ~/go/bin/sup-hcl2-parser

# Assurez-vous que ~/go/bin est dans votre .zshrc
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc
```

## Étape 1 : Initialisation de l'infrastructure locale (Tart)

Nous allons créer un cluster de 4 serveurs Ubuntu.
```bash
# Installation de Tart si nécessaire
brew install cirruslabs/cli/tart

# Récupération de l'image officielle Ubuntu
tart pull ghcr.io/cirruslabs/ubuntu:latest

# Clonage et lancement des 4 noeuds
for i in {1..4}; do 
  tart clone ghcr.io/cirruslabs/ubuntu:latest "node-$i"
  nohup tart run "node-$i" > /dev/null 2>&1 & 
done
```

Attendez 30 secondes, puis récupérez les IPs :
```bash
echo "--- Inventaire du cluster ---"
for i in {1..4}; do echo "node-$i : $(tart ip "node-$i")"; done
```

## Étape 2 : Configuration SSH

Injectez votre clé SSH dans les VMs (le mot de passe par défaut est `admin`) :
```bash
# Remplacez les IPs par celles obtenues ci-dessus
for ip in 192.168.64.XX 192.168.64.XX ...; do
  ssh-copy-id -o StrictHostKeyChecking=no admin@$ip
done
```

## Étape 3 : Définition de la configuration HCL2

Créez un fichier `Supfile.hcl` dans votre dossier de travail. 
**Note importante** : Pour utiliser des variables d'environnement shell dans les commandes, utilisez le double dollar `$$` pour échapper l'interpolation HCL.

```hcl
version = "0.5"

env {
  APP_NAME = "demo-app"
  LOG_PATH = "/var/log/demo.log"
}

network "cluster" {
  hosts = [
    "admin@192.168.64.90", # À mettre à jour avec vos IPs
    "admin@192.168.64.89",
    "admin@192.168.64.88",
    "admin@192.168.64.91"
  ]
}

command "check-system" {
  desc = "Vérifie les ressources système"
  run  = "uname -a && free -h && df -h /"
}

command "prepare-node" {
  desc = "Prépare les répertoires"
  run  = "sudo mkdir -p /opt/$${APP_NAME} && sudo touch $${LOG_PATH} && ls -l $${LOG_PATH}"
}

command "local-audit" {
  desc  = "Audit local"
  local = "echo \"Début de l'audit à : $(date)\" > deployment.audit"
}

target "full-setup" {
  commands = ["local-audit", "check-system", "prepare-node"]
}
```

## Étape 4 : Exécution de la démonstration

Lancez l'orchestration :
```bash
sup -f Supfile.hcl --parser sup-hcl2-parser cluster full-setup
```

## Étape 5 : Nettoyage

```bash
killall tart
for i in {1..4}; do tart delete "node-$i"; done
```
