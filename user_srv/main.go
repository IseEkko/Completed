package main

import (
	"Completed_Server/user_srv/global"
	"Completed_Server/user_srv/handler"
	"Completed_Server/user_srv/initialize"
	"Completed_Server/user_srv/proto"
	"Completed_Server/user_srv/utils"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")
	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	flag.Parse()
	//动态的获取监听端口
	if *Port == 0 {
		*Port, _ = utils.GetFrePort()
	}
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
		GRPC:                           fmt.Sprintf("192.168.6.119:%d", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.SerConfig.Name
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceId
	registration.Port = *Port
	registration.Tags = []string{"imooc", "bobby", "user", "srv"}
	registration.Address = "192.168.6.119"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	/**
	这里进行了阻塞，也就是说这里一直阻塞着，就不会拿到信号量，那么就一直不会退出
	解决的办法：使用异步操作
	*/
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//优雅的退出服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")

}
