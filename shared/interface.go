package shared

import (

	"fmt"
	"net/rpc"
	"github.com/hashicorp/go-plugin"
	
)

var Handshake = plugin.HandshakeConfig {
	ProtocolVersion: 1,
	MagicCookieKey: "SUP_HCL2_PLUGIN",
	MagicCookieValue: "sup-hcl2",
}

var pluginMap = map[string]plugin.Plugin{
	"config_parser": &ConfigParserPlugin{},	

}
type ConfigParser interface{
	ParseFile(filePath string) ([]byte, error) 
}
type ParseFileArgs struct {
	FilePath string
}

type ParseFileReply struct {
	YAMLData []byte
	Error string

}

type ConfigParserRPCClient struct {
	client *rpc.Client
}

func (c *ConfigParserRPCClient) ParseFile(filePath string) ([]byte, error){
	var reply ParseFileReply
	err := c.client.Call("Plugin.Parsefile", &ParseFileArgs{FilePath: filePath}, &reply)
	if err != nil {
		return nil, err
	}
	if reply.Error != ""{
		return nil, fmt.Errorf("%s",reply.Error)
	}
	return reply.YAMLData, nil
}

type ConfigParserRPCServer struct {
	Impl ConfigParser
}


func (s *ConfigParserRPCServer) ParseFile(args *ParseFileArgs, reply *ParseFileReply) error{
	data, err := s.Impl.ParseFile(args.FilePath)
	if err != nil{
		reply.Error = err.Error()
		return nil
	}
	reply.YAMLData = data
	return nil

}

type ConfigParserPlugin struct {
	Impl ConfigParser
}


func (p *ConfigParserPlugin) Server(*plugin.MuxBroker) (interface{}, error)  {
	return &ConfigParserRPCServer{Impl: p.Impl}, nil

}

func (p *ConfigParserPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ConfigParserRPCClient{client: c}, nil

}