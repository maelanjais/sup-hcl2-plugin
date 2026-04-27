package main

import (
	"fmt"
	"sort"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"gopkg.in/yaml.v2"
	"github.com/maelanjais/sup-hcl2-plugin/shared"
)

type HCLConfig struct {
	Version  string       `hcl:"version,optional"`
	Networks []HCLNetwork `hcl:"network,block"`
	Commands []HCLCommand `hcl:"command,block"`
	Targets  []HCLTarget  `hcl:"target,block"`
	Remain   hcl.Body     `hcl:",remain"`
}

type HCLNetwork struct {
	Name      string   `hcl:"name,label"`
	Hosts     []string `hcl:"hosts,optional"`
	Inventory string   `hcl:"inventory,optional"`
	Bastion   string   `hcl:"bastion,optional"`
	Remain    hcl.Body `hcl:",remain"`
}
type HCLCommand struct {
	Name    string      `hcl:"name,label"`
	Desc    string      `hcl:"desc,optional"`
	Run     string      `hcl:"run,optional"`
	Script  string      `hcl:"script,optional"`
	Local   string      `hcl:"local,optional"`
	Stdin   bool        `hcl:"stdin,optional"`
	Once    bool        `hcl:"once,optional"`
	Serial  int         `hcl:"serial,optional"`
	Uploads []HCLUpload `hcl:"upload,block"`
}
type HCLUpload struct {
	Src     string `hcl:"src"`
	Dst     string `hcl:"dst"`
	Exclude string `hcl:"exclude,optional"`
}
type HCLTarget struct {
	Name     string   `hcl:"name,label"`
	Commands []string `hcl:"commands"`
}
// ===== Structures YAML de sortie (compatibles avec sup) =====
type YAMLSupfile struct {
	Version  string                  `yaml:"version"`
	Env      yaml.MapSlice           `yaml:"env,omitempty"`
	Networks map[string]*YAMLNetwork `yaml:"networks"`
	Commands map[string]*YAMLCommand `yaml:"commands"`
	Targets  map[string][]string     `yaml:"targets,omitempty"`
}
type YAMLNetwork struct {
	Env       yaml.MapSlice `yaml:"env,omitempty"`
	Inventory string        `yaml:"inventory,omitempty"`
	Hosts     []string      `yaml:"hosts,omitempty"`
	Bastion   string        `yaml:"bastion,omitempty"`
}
type YAMLCommand struct {
	Desc   string        `yaml:"desc,omitempty"`
	Run    string        `yaml:"run,omitempty"`
	Script string        `yaml:"script,omitempty"`
	Local  string        `yaml:"local,omitempty"`
	Stdin  bool          `yaml:"stdin,omitempty"`
	Once   bool          `yaml:"once,omitempty"`
	Serial int           `yaml:"serial,omitempty"`
	Upload []YAMLUpload  `yaml:"upload,omitempty"`
}
type YAMLUpload struct {
	Src     string `yaml:"src"`
	Dst     string `yaml:"dst"`
	Exclude string `yaml:"exclude,omitempty"`
}

type HCL2Parser struct{}

func (p *HCL2Parser) ParseFile(filePath string) ([]byte, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(filePath)
	if diags.HasErrors(){
		return nil, fmt.Errorf("HCL parse error: %s", diags.Error())
	}
	var config HCLConfig
	diags = gohcl.DecodeBody(file.Body,nil, &config)
	if diags.HasErrors(){
		return nil, fmt.Errorf("HCL has decode error: %s",diags.Error())
	}
	globalEnv := extractEnvFromBody(config.Remain)

	out:= &YAMLSupfile{
		Version: config.Version,
		Env: globalEnv,
		Networks: make(map[string]*YAMLNetwork),
		Commands: make(map[string]*YAMLCommand),
		Targets: make(map[string][]string),
	}
	if out.Version == ""{
		out.Version ="0.5"
	}
	// Networks

	for _, n := range config.Networks{
		net := &YAMLNetwork{
			Hosts: n.Hosts,
			Inventory: n.Inventory,
			Bastion: n.Bastion,
			Env: extractEnvFromBody(n.Remain),

		}
		out.Networks[n.Name] = net
	}

	// Commands
	for _, c := range config.Commands{
		cmd := &YAMLCommand{
			Desc: c.Desc,
			Run: c.Run,
			Script: c.Script,
			Local: c.Local,
			Stdin: c.Stdin,
			Once: c.Once,
			Serial: c.Serial,


		}
		for _, u := range c.Uploads{
			cmd.Upload = append(cmd.Upload, YAMLUpload{
				Src: u.Src,
				Dst:u.Dst,
				Exclude: 
				u.Exclude,
			})
		}
		out.Commands[c.Name] = cmd
	
	}

	// Targets
	for _, t := range config.Targets{
		out.Targets[t.Name] = t.Commands
	}
	return yaml.Marshal(out)
}

// extractEnvFromBody extrait les variables d'environnement
// depuis un bloc "env { KEY = VALUE }" dans un hcl.Body restant.

func extractEnvFromBody(body hcl.Body) yaml.MapSlice{
	if body == nil{
		return nil
	}
	content, _, _ := body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "env"},
		},
	})

	if content == nil || len(content.Blocks) == 0{
		return nil
	}
	attrs, _ := content.Blocks[0].Body.JustAttributes()
	if len(attrs) == 0{
		return nil
	}

	keys := make([]string, 0, len(attrs))
	for k :=range attrs{
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var envSlice yaml.MapSlice
	for _, k := range keys{
		val, diags := attrs[k].Expr.Value((nil))
		if diags.HasErrors(){
			continue
		}
		envSlice = append(envSlice, yaml.MapItem{
			Key: k,
			Value: val.AsString(),
		})
	}
	return envSlice
}


func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"config_parser": &shared.ConfigParserPlugin{Impl: &HCL2Parser{}},
		},
	})
}