# Mise à jour suggérée pour le README de ton fork de `sup`

*(Copie-colle cette section dans le fichier `README.md` de ton dossier `sup` pour expliquer ta nouvelle fonctionnalité !)*

---

## 🌟 NOUVEAU : Support natif du HCL2 (via Plugin)

En plus des classiques `Supfile` en YAML, **ce fork de `sup` supporte désormais la syntaxe HashiCorp Configuration Language (HCL2) !**

Le support HCL2 permet une meilleure lisibilité, l'utilisation de blocs explicites, et s'intègre parfaitement avec vos habitudes Terraform/Packer.

### Comment l'utiliser ?

1. **Installer le plugin parseur** : Vous devez avoir l'exécutable `sup-hcl2-parser` dans votre système.
   [👉 Voir le dépôt du plugin HCL2](https://github.com/maelanjais/sup-hcl2-plugin)

2. **Écrire un `Supfile.hcl`** :
```hcl
version = "0.5"

network "production" {
  hosts = ["admin@api1.domain.com", "admin@api2.domain.com"]
}

command "deploy" {
  desc = "Update the API"
  run  = "docker pull api:latest && docker restart api"
}
```

3. **Lancer Sup** : 
```bash
sup -f Supfile.hcl production deploy
```
*(Sup détectera automatiquement l'extension `.hcl` et appellera le plugin en arrière-plan pour traiter le fichier !)*
