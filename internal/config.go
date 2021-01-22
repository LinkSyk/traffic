package internal

type TrafficConfig struct {
	TcpUpstreams []*TcpListenerConfig `json:"tcpUpstreams" yaml:"tcpUpstreams"`
}

type TcpListenerConfig struct {
	Name string `json:"name" yaml:"name"`
	Listen string `json:"listen" yaml:"listen"`
	LbAlg string `json:"lbAlg" yaml:"lbAlg"`
	UpstreamsConfig []NodeConfig `json:"upstreams" yaml:"upstreams"`
}

type NodeConfig struct {
	Name string `json:"name" yaml:"name"`
	Addr string `json:"addr" yaml:"addr"`
	Weight float64 `json:"weight" yaml:"weight"`
}
