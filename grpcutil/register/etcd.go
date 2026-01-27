package register

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"strings"
	"time"

	etcd "go.etcd.io/etcd/client/v3"

	registry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	ggrpc "google.golang.org/grpc"
)

type EtcdConf struct {
	Env       string   `yaml:"env"`
	Endpoints []string `yaml:"endpoints"`
	TlsCert   string   `yaml:"tls_cert"`
	TlsKey    string   `yaml:"tls_key"`
	TlsCa     string   `yaml:"tls_ca"`
	Token     string   `yaml:"token"`
}

func (e *EtcdConf) GetTlsConfig() *tls.Config {
	if e.TlsCert == "" || e.TlsKey == "" || e.TlsCa == "" {
		log.Panicf("[etcd] etcd tls cert content is empty,lenght for cert=%d, key=%d, ca=%d,", len(e.TlsCert), len(e.TlsKey), len(e.TlsCa))
	}

	ct, err := tls.X509KeyPair([]byte(e.TlsCert), []byte(e.TlsKey))
	if err != nil {
		log.Panicf("[etcd] parses etcd tls cert error:%s", err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(e.TlsCa))

	return &tls.Config{
		Certificates: []tls.Certificate{ct},
		RootCAs:      pool,
	}
}

func NewEtcdRegistry(conf *EtcdConf) *registry.Registry {
	if len(conf.Endpoints) == 0 {
		log.Panicf("[ectd] Failed to ectd endpoints empty")
	}

	log.Println("[etcd] connect:", strings.Join(conf.Endpoints, ","))

	cfg := etcd.Config{
		Endpoints:            conf.Endpoints,
		DialTimeout:          2 * time.Second,
		DialKeepAliveTime:    5 * time.Second, // 每5秒ping一次服务器
		DialKeepAliveTimeout: time.Second,     // 1秒没有返回则代表故障
		DialOptions:          []ggrpc.DialOption{ggrpc.WithBlock()},
	}

	if conf.TlsCert != "" {
		cfg.TLS = conf.GetTlsConfig()
	}

	cli, err := etcd.New(cfg)
	if err != nil {
		log.Panicf("[etcd] Failed to etcd:%s", err)
	}

	if conf.Env == "" {
		log.Panicf("[etcd] Failed to etcd config must set env")
	}

	return registry.New(cli, registry.Namespace(fmt.Sprintf("/services/%s", conf.Env)))
}
