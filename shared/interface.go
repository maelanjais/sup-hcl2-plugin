package shared

import (
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Handshake partagé entre le host (sup) et le plugin.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SUP_HCL2_PLUGIN",
	MagicCookieValue: "sup-hcl2",
}

// PluginMap est la map des plugins disponibles.
var PluginMap = map[string]plugin.Plugin{
	"config_parser": &ConfigParserPlugin{},
}

// ConfigParser est l'interface exposée par le plugin.
// ParseFile lit un fichier HCL2 et retourne des bytes YAML
// compatibles avec sup.NewSupfile().
type ConfigParser interface {
	ParseFile(filePath string) ([]byte, error)
}

// --- Arguments et réponse RPC ---

type ParseFileArgs struct {
	FilePath string
}

type ParseFileReply struct {
	YAMLData []byte
	Error    string
}

// --- Client RPC (côté host/sup) ---

type ConfigParserRPCClient struct {
	client *rpc.Client
}

func (c *ConfigParserRPCClient) ParseFile(filePath string) ([]byte, error) {
	var reply ParseFileReply
	err := c.client.Call("Plugin.ParseFile", &ParseFileArgs{FilePath: filePath}, &reply)
	if err != nil {
		return nil, err
	}
	if reply.Error != "" {
		return nil, fmt.Errorf("%s", reply.Error)
	}
	return reply.YAMLData, nil
}

// --- Serveur RPC (côté plugin) ---

type ConfigParserRPCServer struct {
	Impl ConfigParser
}

func (s *ConfigParserRPCServer) ParseFile(args *ParseFileArgs, reply *ParseFileReply) error {
	data, err := s.Impl.ParseFile(args.FilePath)
	if err != nil {
		reply.Error = err.Error()
		return nil
	}
	reply.YAMLData = data
	return nil
}

// --- Plugin wrapper (implémente plugin.Plugin) ---

type ConfigParserPlugin struct {
	Impl ConfigParser
}

func (p *ConfigParserPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ConfigParserRPCServer{Impl: p.Impl}, nil
}

func (p *ConfigParserPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ConfigParserRPCClient{client: c}, nil
}
