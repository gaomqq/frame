package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strconv"
)

var client config_client.IConfigClient

func GetClient(ip, port string) error {

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, uint64(portInt), constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	return err
}

func GetConfig(dataID, group string) (string, error) {
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "user",
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		return "", err
	}

	return content, nil
}

// TODO:完成对mysql的监听
func ListenConfig(group, dataID string) error {
	return client.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})
}
