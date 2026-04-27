# Supfile.hcl — Configuration HCL2 pour sup

version = "0.5"

# Variables d'environnement globales
env {
  NAME  = "api"
  IMAGE = "example/api"
}

# === Réseaux ===

network "local" {
  hosts = ["localhost"]
}

network "staging" {
  hosts = ["stg1.example.com"]

  env {
    ENVIRONMENT = "staging"
  }
}

network "production" {
  hosts   = ["api1.example.com", "api2.example.com"]
  bastion = "bastion.example.com"

  env {
    ENVIRONMENT = "production"
  }
}

# === Commandes ===

command "echo" {
  desc = "Print some env vars"
  run  = "echo $NAME $IMAGE $SUP_NETWORK"
}

command "date" {
  desc = "Print OS name and current date/time"
  run  = "uname -a; date"
}

command "build" {
  desc = "Build Docker image and push to registry"
  run  = "sudo docker build -t $IMAGE:latest . && sudo docker push $IMAGE:latest"
  once = true
}

command "pull" {
  desc = "Pull latest Docker image"
  run  = "sudo docker pull $IMAGE:latest"
}

command "restart" {
  desc   = "Restart Docker container"
  run    = "sudo docker restart $NAME"
  serial = 2
}

command "bash" {
  desc  = "Interactive Bash on all hosts"
  stdin = true
  run   = "bash"
}

command "upload-dist" {
  desc = "Upload dist files to all hosts"

  upload {
    src = "./dist"
    dst = "/tmp/"
  }
}

command "prepare" {
  desc  = "Prepare to upload"
  local = "npm run build"
}

# === Targets (alias pour plusieurs commandes) ===

target "deploy" {
  commands = ["build", "pull", "restart"]
}

target "all" {
  commands = ["echo", "date"]
}
