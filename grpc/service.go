package grpc

import (
	"encoding/json"
	"fmt"
	"github.com/gaomqq/frame/config"
	"github.com/gaomqq/frame/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Config struct {
	App struct {
		Ip   string `yaml:"ip" json:"ip"`
		Port string `yaml:"port" json:"port"`
	} `yaml:"app" json:"app"`
}

func getConfig(nacosGroup, serviceName string) (*Config, error) {
	configInfo, err := config.GetConfig(nacosGroup, serviceName)
	if err != nil {
		return nil, err
	}
	cnf := new(Config)
	err = json.Unmarshal([]byte(configInfo), cnf)
	if err != nil {
		return nil, err
	}
	return cnf, nil
}

func RegisterGRPC(nacosGroup, serviceName string, register func(s *grpc.Server)) error {
	cof, err := getConfig(nacosGroup, serviceName)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "0.0.0.0", "8081"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	s := grpc.NewServer()
	//反射接口支持查询
	reflection.Register(s)
	//支持健康检查
	healthpb.RegisterHealthServer(s, health.NewServer())
	fmt.Println(cof.App.Ip, cof.App.Port)

	err = consul.ServiceRegister(nacosGroup, serviceName, cof.App.Ip, serviceName)
	fmt.Println(cof.App.Ip)
	fmt.Println(cof.App.Port)
	if err != nil {
		return err
	}

	register(s)
	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
