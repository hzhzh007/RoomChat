// load config file
// here used hot reload, but it may not be used
package main

import (
	"flag"
	log "github.com/hzhzh007/RoomChat/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	//"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ConnectorConf struct {
	Hostname  string `yaml:"Hostname"`
	ServeAddr string `yaml:"ServeAddr"`
	RpcAddr   string `yaml:"RpcAddr"`
}

type Config struct {
	ServeAddr  string          `yaml:"ServeAddr"`
	Connectors []ConnectorConf `yaml:"Connectors"`
	Log        log.LogConfig   `yaml:"log"`
}

func NewDefaultConfig() *Config {
	return &Config{
		ServeAddr: "localhost:8080",
		Connectors: []ConnectorConf{
			{Hostname: "1", ServeAddr: "127.0.0.1:8888", RpcAddr: "127.0.0.1:8888"},
			{Hostname: "1", ServeAddr: "127.0.0.1:8889", RpcAddr: "127.0.0.1:8889"},
		},
		Log: log.LogConfig{Module: "connector",
			FileName: "log.log",
			Level:    1,
			Format:   "%{pid} %{time} %{module} %{shortfile} %{message}",
		},
	}
}

var (
	Conf       *Config
	confFile   string
	configLock = new(sync.RWMutex)
)

func init() {
	flag.StringVar(&confFile, "c", "request.yaml", "set the connector conf file")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			loadConfig(false)
			log.Info("Reloaded")
		}
	}()
}

func GetConfig() *Config {
	if Conf == nil {
		loadConfig(true)
	}
	configLock.RLock()
	defer configLock.RUnlock()
	return Conf
}

func loadConfig(fail bool) error {
	temp := NewDefaultConfig()
	if Conf == nil {
		configLock.Lock()
		Conf = temp
		configLock.Unlock()
	}
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Error("open config: ", err)
		if fail {
			os.Exit(1)
		}
		return err
	}
	if err = yaml.Unmarshal(data, temp); err != nil {
		log.Error("yaml unmarshal error", err)
		if fail {
			os.Exit(1)
		}
		return err
	}
	log.Info("load config ok", *temp)
	configLock.Lock()
	Conf = temp
	configLock.Unlock()
	return nil
}