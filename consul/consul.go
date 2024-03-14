package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gaomqq/frame/config"
	"github.com/gaomqq/frame/redis"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"net"
	"strconv"
	"time"
)

const CONSUL_KEY = "consul:node:index"

type ConsulConfig struct {
	Consul struct {
		Ip   string `yaml:"ip" json:"ip"`
		Port string `yaml:"port" json:"port"`
	} `yaml:"consul" json:"consul"`
}

func getConfig(nacosGroup, serviceName string) (*ConsulConfig, error) {
	cnf, err := config.GetConfig(nacosGroup, serviceName)
	if err != nil {
		return nil, err
	}

	consulCnf := new(ConsulConfig)
	err = json.Unmarshal([]byte(cnf), consulCnf)
	if err != nil {
		return nil, err
	}

	return consulCnf, err
}

func getIndex(ctx context.Context, serviceName string, indexLen int) (int, error) {
	exist, err := redis.ExistKey(ctx, serviceName, CONSUL_KEY)
	if err != nil {
		return 0, err
	}

	if exist {
		indexStr, err := redis.GetByKey(ctx, serviceName, CONSUL_KEY)
		if err != nil {
			return 0, err
		}
		index, err := strconv.Atoi(indexStr)
		newIndex := index + 1

		if newIndex >= indexLen {
			newIndex = 0
		}
		err = redis.SetKey(ctx, serviceName, CONSUL_KEY, newIndex, time.Duration(0))
		if err != nil {
			return 0, err
		}

		return index, nil
	}

	err = redis.SetKey(ctx, serviceName, "consul:node:index", 0, time.Duration(0))
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func AgentHealthService(ctx context.Context, serviceName string) (string, error) {
	client, err := capi.NewClient(capi.DefaultConfig())
	if err != nil {
		return "", err
	}
	sr, infos, err := client.Agent().AgentHealthServiceByName(serviceName)
	if err != nil {
		return "", err
	}
	if sr != "passing" {
		return "", fmt.Errorf("is not have health service")
	}

	index, err := getIndex(ctx, serviceName, len(infos))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v:%v", infos[index].Service.Address, infos[index].Service.Port), nil
}

func ServiceRegister(nacosGroup, serviceName string, address string, port string) error {

	config:= capi.DefaultConfig()
	config.Address = "10.2.171.125:8500"
	client,_:=capi.NewClient(config)
	return client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    "user",
		Tags:    []string{"GRPC"},
		Port:    8081,
		Address: GetIp()[0],
		Check: &capi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%v:%v", GetIp()[0], "8081"),
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func GetIp() (ip []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, addr := range addrs {
		ipNet, isVailIpNet := addr.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = append(ip, ipNet.IP.String())
			}
		}

	}
	return ip
}
