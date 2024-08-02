package grpc

type PortConfig struct {
	grpcPort      string
	httpProxyPort string
	wsProxyPort   string
}

func NewPortConfig(grpcPort string, httpProxyPort string, wsProxyPort string) *PortConfig {
	return &PortConfig{
		grpcPort:      grpcPort,
		httpProxyPort: httpProxyPort,
		wsProxyPort:   wsProxyPort,
	}
}

func (p *PortConfig) GRPCPort() string {
	return p.grpcPort
}
