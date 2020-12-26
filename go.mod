module github.com/hero1s/golib

go 1.14

require (
	git.moumentei.com/plat_go/golib v0.0.0-20201225013707-2916df43f2a2
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.623
	github.com/aliyun/aliyun-oss-go-sdk v2.1.4+incompatible
	github.com/bsm/redislock v0.5.0
	github.com/bsm/sarama-cluster v2.1.15+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/denverdino/aliyungo v0.0.0-20201110010600-f7c7d0e1d041
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gansidui/geohash v0.0.0-20141019080235-ebe5ba447f34
	github.com/gansidui/nearest v0.0.0-20141019122829-a5d0cde6ef14
	github.com/garyburd/redigo v1.6.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/golang-lru v0.5.4
	github.com/howeyc/fsnotify v0.9.0
	github.com/iGoogle-ink/gopay v1.3.9
	github.com/lib/pq v1.7.0
	github.com/logrusorgru/aurora v0.0.0-20181002194514-a7b3b318ed4e
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/kafka/v2 v2.9.1
	github.com/micro/go-plugins/broker/mqtt/v2 v2.9.1
	github.com/micro/go-plugins/broker/nsq/v2 v2.9.1
	github.com/micro/go-plugins/broker/redis/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/breaker/hystrix/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1
	github.com/nacos-group/nacos-sdk-go v1.0.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/opentracing/opentracing-go v1.1.0
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.0
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/swaggo/swag v1.6.5
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/zheng-ji/goSnowFlake v0.0.0-20180906112711-fc763800eec9
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e
	google.golang.org/api v0.35.0
	google.golang.org/grpc v1.31.1
	gopkg.in/Shopify/sarama.v1 v1.20.1
	gopkg.in/fatih/set.v0 v0.2.1
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

replace github.com/iGoogle-ink/gopay v1.3.9 => github.com/hero1s/gopay v1.4.3

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/swaggo/swag v1.6.5 => github.com/nicle-lin/swag v1.6.6-0.20201119025013-fa0209863a28
