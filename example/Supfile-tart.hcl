version = "0.5"

# Réseau avec tes 4 VMs Ubuntu
network "tart-cluster" {
  hosts = [
    "admin@192.168.64.84",
    "admin@192.168.64.87",
    "admin@192.168.64.86",
    "admin@192.168.64.85"
  ]
}

# Affiche les infos système
command "status" {
  desc = "Affiche les infos système et l'uptime de la VM"
  run  = "uname -a && uptime"
}

# Crée un fichier de test
command "ping-test" {
  desc = "Crée un petit fichier pour prouver le passage de sup"
  run  = "echo 'Hello from HCL2 Plugin!' > /tmp/sup-was-here.txt && cat /tmp/sup-was-here.txt"
}
