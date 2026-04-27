package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	hcplugin "github.com/hashicorp/go-plugin"

	"github.com/maelanjais/sup-hcl2-plugin/shared"
)

func main() {
	log.SetOutput(io.Discard)

	client := hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("./sup-hcl2-parser"),
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection error: %v\n", err)
		os.Exit(1)
	}

	raw, err := rpcClient.Dispense("config_parser")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Dispense error: %v\n", err)
		os.Exit(1)
	}

	parser := raw.(shared.ConfigParser)
	yamlData, err := parser.ParseFile("example/Supfile.hcl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(yamlData))
}
