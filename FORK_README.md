# Sup (Edition HCL2)

Ce dépôt est un fork officiel de [pressly/sup](https://github.com/pressly/sup) intégrant le support natif des fichiers de configuration au format **HCL2**.

Sup est un outil de déploiement minimaliste conçu pour l'exécution de commandes sur des flottes de serveurs via SSH. Cette version étend les capacités de Sup en permettant l'utilisation du langage de configuration de Terraform pour définir vos réseaux, commandes et cibles.

## Caractéristiques principales

- **Parsing HCL2** : Utilisation de la syntaxe HCL2 pour une meilleure lisibilité et une gestion simplifiée des variables.
- **Architecture modulaire** : L'intégration repose sur un système de plugins RPC, préservant la légèreté du binaire original.
- **Compatibilité descendante** : Les fichiers YAML traditionnels restent parfaitement supportés.

## Installation

### 1. Compilation de Sup
Clonez ce dépôt et compilez le binaire principal :
```bash
go build -o sup ./cmd/sup/
```

### 2. Installation du parseur HCL2
Pour traiter les fichiers `.hcl`, Sup nécessite le plugin externe `sup-hcl2-parser`. Vous pouvez le compiler depuis le dépôt dédié : [sup-hcl2-plugin](https://github.com/maelanjais/sup-hcl2-plugin).

## Guide d'utilisation

Créez un fichier nommé `Supfile.hcl` à la racine de votre projet :

```hcl
version = "0.5"

network "production" {
  hosts = ["admin@srv1.example.com", "admin@srv2.example.com"]
}

command "update" {
  desc = "Mise à jour de l'application"
  run  = "git pull origin main && make restart"
}
```

Exécutez ensuite votre déploiement :
```bash
./sup -f Supfile.hcl production update
```

## Contribution et Développement

Ce fork a été développé pour répondre aux besoins de modularité et de clarté dans les fichiers de déploiement. Les contributions sont les bienvenues via des Pull Requests.
