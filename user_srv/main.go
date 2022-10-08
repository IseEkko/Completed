package main

import (
	"Completed_Server/user_srv/global"
	"Completed_Server/user_srv/handler"
	"Completed_Server/user_srv/initialize"
	"Completed_Server/user_srv/proto"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")
	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	flag.Parse()
	zap.S().Info("ip:", *IP)
	zap.S().Info("port:", *Port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.SerConfig.CosulConfigInfo.Host,
		global.SerConfig.CosulConfigInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	check := &api.AgentServiceCheck{
		GRPC:                           "192.168.6.119:50051",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.SerConfig.Name
	registration.ID = global.SerConfig.Name
	registration.Port = *Port
	registration.Tags = []string{"imooc", "bobby", "user", "srv"}
	registration.Address = "192.168.6.119"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
