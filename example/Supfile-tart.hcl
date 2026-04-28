version = "0.5"

env {
  APP_NAME = "demo-app"
  LOG_PATH = "/var/log/demo.log"
}

network "tart-cluster" {
  hosts = [
    "admin@192.168.64.92",
    "admin@192.168.64.93",
    "admin@192.168.64.94",
    "admin@192.168.64.95"
  ]
}

command "check-system" {
  desc = "Vérifie les ressources système sur tous les noeuds"
  run  = "uname -a && uptime && free -h && df -h /"
}

command "prepare-node" {
  desc = "Prépare les répertoires et fichiers de log"
  run  = "sudo mkdir -p /opt/$${APP_NAME} && sudo touch $${LOG_PATH} && ls -l $${LOG_PATH}"
}

command "local-audit" {
  desc  = "Génère un rapport local de début de déploiement"
  local = "echo \"Début de l'audit à : $(date)\" > deployment.audit"
}

target "full-setup" {
  commands = ["local-audit", "check-system", "prepare-node"]
}

