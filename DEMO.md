# HCL2 Plugin for Sup — Live Demo 🚀

Voici un scénario de bout en bout pour démontrer la puissance du plugin HCL2 avec de vraies machines virtuelles.

## Prérequis
- Avoir compilé `sup` avec le support du plugin.
- Avoir compilé `sup-hcl2-parser`.
- Avoir `tart` installé sur macOS.

## 1. Démarrer le cluster de VMs (Tart)

Lancez 3 machines Ubuntu en arrière-plan :

```bash
# Cloner les machines
tart clone ghcr.io/cirruslabs/ubuntu:latest vm1
tart clone ghcr.io/cirruslabs/ubuntu:latest vm2
tart clone ghcr.io/cirruslabs/ubuntu:latest vm3

# Démarrer en tâche de fond
nohup tart run vm1 > /dev/null 2>&1 &
nohup tart run vm2 > /dev/null 2>&1 &
nohup tart run vm3 > /dev/null 2>&1 &
```

Récupérez leurs adresses IP :
```bash
tart ip vm1 && tart ip vm2 && tart ip vm3
```

## 2. Configurer les clés SSH

Pour que `sup` puisse s'y connecter sans mot de passe, ajoutez votre clé SSH (le mot de passe par défaut de ces VMs est `admin`) :

```bash
# Remplacez les IP par celles obtenues à l'étape 1
for ip in 192.168.64.84 192.168.64.85 192.168.64.86; do
  sshpass -p admin ssh-copy-id -o StrictHostKeyChecking=no admin@$ip
done
```

## 3. Le fichier HCL2 (`Supfile-demo.hcl`)

Créez ce fichier, c'est lui qui orchestre le tout (en remplaçant les IP) :

```hcl
version = "0.5"

network "demo-cluster" {
  hosts = [
    "admin@192.168.64.84",
    "admin@192.168.64.85",
    "admin@192.168.64.86"
  ]
}

command "status" {
  desc = "Affiche l'OS et la charge système"
  run  = "uname -a && uptime"
}

command "ping-test" {
  desc = "Prouve que sup est passé par là"
  run  = "echo 'Hello from HCL2 Plugin!' > /tmp/sup-was-here.txt && cat /tmp/sup-was-here.txt"
}

target "deploy" {
  commands = ["status", "ping-test"]
}
```

## 4. Lancement de la Démo !

Exécutez `sup` en pointant vers le plugin et le fichier `.hcl` :

```bash
./sup -f Supfile-demo.hcl --parser ./sup-hcl2-parser demo-cluster deploy
```

Vous verrez les 3 serveurs répondre en parallèle ! 🎉
