# qbtool

糗百go工具
功能点:
    1:db结构体生成go代码
    2:db结构体差异对比生成sql
    3:swag 文档生成工具


安装
    make install

参数说明
qbtool help
qbtool help command

示例:
生成文档
qbtool swag init -d=api -g=main.go -o=api/swagger

db结构体生成go代码
qbtool db reverse -source="root:e23456@tcp(172.16.3.21:8306)/cherry?charset=utf8" -path="./" -single=false

db对比源
qbtool db diff -source="root:e23456@tcp(172.16.3.21:8306)/cherry?charset=utf8" -target="root:e23456@tcp(172.16.3.21:8306)/test1?charset=utf8" -path="./diff.sql"



利用模板生成业务代码样例

const rpcTemplateText = `syntax = "proto3";

package {{.package}};

message Request {
  string ping = 1;
}

message Response {
  string pong = 1;
}

service {{.serviceName}} {
  rpc Ping(Request) returns(Response);
}
`
err = util.With("t").Parse(text).SaveTo(map[string]string{
    "package":     serviceName.UnTitle(),
    "serviceName": serviceName.Title(),
}, out, false)







