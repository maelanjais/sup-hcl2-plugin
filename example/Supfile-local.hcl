version = "0.5"

network "my-mac" {
  # On dit à sup de cibler ta propre machine
  hosts = ["localhost"]
}

command "hello-mac" {
  desc  = "Affiche un message de succès depuis ton Mac"
  # 'local' dit à sup d'exécuter la commande directement sur ta machine (sans SSH)
  local = "echo '=====================================' && echo '🚀 SUCCÈS : Le parseur HCL2 de Sup a parfaitement fonctionné !' && echo '=====================================' && uname -a"
}

command "list-files" {
  desc  = "Liste les fichiers de ton projet"
  local = "ls -la"
}

target "demo" {
  commands = ["hello-mac", "list-files"]
}
