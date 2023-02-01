package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

// ProviderConfig is config struct
type ProviderConfig struct {
	NodeURL                          string
	NodeChainID                      string
	CPUPrice                         string
	MemoryPrice                      string
	StoragePrice                     string
	SecretKey                        string
	OrderFactory                     string
	ProviderContract                 string
	ProviderFactoryContract          string
	BidTimeOut                       int64
	NameSpace                        string
	K8sConfigPath                    string
	HostNameServiceListenAddr        string
	Cert                             string
	GatewayListenAddress             string
	ProviderAddress                  string
	DeploymentIngressDomain          string
	DeploymentIngressExposeLBHosts   bool
	DeploymentIngressStaticHosts     bool
	DeploymentNetworkPoliciesEnabled bool
	ClusterPublicHostname            string
	OvercommitPercentCPU             float64
	OvercommitPercentMemory          float64
	OvercommitPercentStorage         float64
	DeploymentRuntimeClass           string
	DockerImagePullSecretsName       string
	HostPruneInterval                int
	HostWebRefreshInterval           int
	HostRetryDelay                   int
}

// LoadConfig get config from config.json
func LoadConfig(c *viper.Viper) *ProviderConfig {
	//c = viper.New()
	//fmt.Println("c1",c)
	c.SetConfigName("config")
	c.AddConfigPath(".")
	c.SetConfigType("json")
	c.SetDefault("OrderFactory", "")
	if err := c.ReadInConfig(); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	log.Println("load config success", c.GetString("OrderFactory"))

	return ConvertConfig(c)
}

// WatchConfig is function watching config change
func WatchConfig(c *viper.Viper) error {
	if err := LoadConfig(c); err == nil {
		return errors.New("unable load config")
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.WatchConfig()
	watch := func(e fsnotify.Event) {
		log.Printf("Config file is changed: %s\n", e.String())
		cancel()
	}
	c.OnConfigChange(watch)
	<-ctx.Done()
	return nil
}

// ConvertConfig is convent config function
func ConvertConfig(c *viper.Viper) *ProviderConfig {
	pConfig := ProviderConfig{}
	fmt.Println(c.GetString("OrderFactory"))
	if err := c.Unmarshal(&pConfig); err != nil {
		log.Println("convertConfig error", err.Error())
		return nil
	}
	return &pConfig
}

// ConvertConfigPtr is convent config struct function
func ConvertConfigPtr(c *viper.Viper, pConfig *ProviderConfig) error {
	//pConfig := ProviderConfig{}
	//fmt.Println(c.GetString("OrderFactory"))
	if err := c.Unmarshal(pConfig); err != nil {
		log.Println("convertConfig error", err.Error())
		return err
	}
	return nil
}
