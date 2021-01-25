package internal

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type TrafficConfig struct {
	TcpUpstreams []*TcpListenerConfig `json:"tcpUpstreams" yaml:"tcpUpstreams"`
}

type TcpListenerConfig struct {
	Name   string       `json:"name" yaml:"name"`
	Listen string       `json:"listen" yaml:"listen"`
	LbAlg  string       `json:"lbAlg" yaml:"lbAlg"`
	Nodes  []NodeConfig `json:"upstreams" yaml:"upstreams"`
}

type NodeConfig struct {
	Name   string  `json:"name" yaml:"name"`
	Addr   string  `json:"addr" yaml:"addr"`
	Weight float64 `json:"weight" yaml:"weight"`
}

// 解析配置、验证配置

func CheckConfig(cfg *TrafficConfig) error {
	return nil
}

func ParserConfig(configReader io.Reader) (*TrafficConfig, error) {
	data, err := ioutil.ReadAll(configReader)
	if err != nil {
		return nil, err
	}

	tc := new(TrafficConfig)
	if err := yaml.Unmarshal(data, tc); err != nil {
		return nil, err
	}

	return tc, nil
}
