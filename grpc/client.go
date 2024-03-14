package grpc

import (
	"context"
	"fmt"
	"github.com/gaomqq/frame/consul"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client2(ctx context.Context, toService string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("consul://10.2.171.125:8500/"+toService+"?wait=14s", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`))

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return grpc.Dial("consul://10.2.171.125:8500/"+toService+"?wait=14s", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`))

}

func Client(ctx context.Context, toService string) (*grpc.ClientConn, error) {
	conn, err := consul.AgentHealthService(ctx, toService)
	if err != nil {
		return nil, err
	}
	fmt.Println(conn)
	return grpc.Dial(conn, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
