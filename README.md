# golibs

#### 介绍
golang 业务基础库

#### 规范
#####添加目录必须同步说明。大的模块内部编写readme.md文档
#####1 不得引入写死配置的功能模块
#####2 尽量不编写单例模块，以对象或参数实现复用
#####3 引入的业务基础库必须测试可用，复杂模块保留测试用例
#####4 文件夹,文件名小写 
#####5 以tag版本的形式引用,尽量在项目中测试通过再入库
#####6 同类型的代码归集到相应的目录

#### 目录
##### cache      	与缓存相关的操作(限制,随机ID,分布式定时器,缓存便捷操作)

##### conf 			配置文件读取

##### connsvr 	 长连接服务框架  

##### constant     与业务逻辑无关的常量
##### db         		数据库的操作封装

##### Event   		事件服务

##### grpc       	  grpc相关工具
##### helpers       常用小工具

##### i18n  			错误码相关 

##### log        		日志库

##### micro  		 go-micro相关封装

##### proxy      	 代理转发
##### stringutils  字符串操作

##### task 			定时任务

##### third_sdk    第三方SDK接入
##### tools      	  业务工具集
##### utils      		核心工具

##### validation   参数校验库

##### watch 		  文件监控

##### web         	  web工具集,micro服务快捷启动

#### 安装go micro v2 版本，限定使用v2版本
##### install protoc-gen-go
go get github.com/golang/protobuf/{proto,protoc-gen-go}
##### install protoc-gen-micro
go get github.com/micro/micro/v2/cmd/protoc-gen-micro@master

#### 安装jaeger 测试
docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p5775:5775/udp -p6831:6831/udp -p6832:6832/udp  -p5778:5778 -p16686:16686 -p14268:14268 -p9411:9411 jaegertracing/all-in-one:latest
http://localhost:16686