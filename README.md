# Completed项目实现
   从20220928开始制作项目，项目背景：
   随着人们生活水平的提高，对电商的安全性和稳定性的追求日益增长，目前传统架构的电商系统，稳定性查，维护成本高等缺点，对整个系统和用户体验来说都是极差的。Grpc框架的微服务电商系统，是一个基于google开源框架GRPC进行实现的web系统，具有高传输性、安全性、易维护性等优点。使用GRPC解决了传统使用Json格式进行传输的低效性。同时GRPC底层为http2实现的传输协议，传输效率也得到了极大的提高
## 用户服务
```txt
├── go.mod
├── go.sum
└── user_srv
    ├── Tests
    │   └── main.go
    ├── global
    │   └── global.go
    ├── handler
    │   └── user.go
    ├── main.go
    ├── model
    │   ├── main
    │   │   └── main.go
    │   └── user.go
    └── proto
        ├── user.pb.go
        └── user.proto
 ```

  * 获取用户列表
  * 通过手机号码进行查询
  * 通过id进行查询
  * 用户个人中心的信息更改
  * 密码校验
