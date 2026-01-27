package register

import (
	"fmt"
	"log"

	registry "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	consul "github.com/hashicorp/consul/api"
)

type ConsulConf struct {
	Env     string `yaml:"env"`
	Address string `yaml:"address"`
	Scheme  string `yaml:"scheme"`
	TlsCert string `yaml:"tls_cert"`
	TlsKey  string `yaml:"tls_key"`
	TlsCa   string `yaml:"tls_ca"`
	Token   string `yaml:"token"`
}

func NewConsulRegistry(conf *ConsulConf) *registry.Registry {
	if len(conf.Address) == 0 {
		log.Panicf("[consul] Failed to consul endpoints empty")
	}

	log.Println("[consul] connect:", conf.Address)

	cfg := consul.Config{
		Address:   conf.Address,
		Scheme:    conf.Scheme,
		Namespace: fmt.Sprintf("/services/%s", conf.Env),
	}

	cli, err := consul.NewClient(&cfg)
	if err != nil {
		log.Panicf("[consul] Failed to consul:%s", err)
	}

	if conf.Env == "" {
		log.Panicf("[consul] Failed to consul config must set env")
	}

	return registry.New(cli)
}
