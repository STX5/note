# Note📝
## 简介
本项目是一个简单的Golang微服务架构web应用，总共有三个微服务：api、note、user。其中 api 服务是对外暴露的 HTTP 接口，采用 Hertz 框架。user 服务与 note 服务为内部的微服务。微服务之间通过 kitex 框架的 RPC 功能进行交互。并且接入Opentelemetry和Jaeger进行观测与链路追踪。
``` bash
├── cmd
│   ├── api             # 对外暴露的 HTTP 接口服务
│   │   └── main.go     # api 服务的入口文件
        └── ...
│   ├── note            # 内部的笔记微服务
│   │   └── main.go     # note 服务的入口文件
        └── ...
│   ├── user            # 内部的用户微服务
│   │   └── main.go     # user 服务的入口文件
        └── ...

```
| ServiceName | Usage                     | Path     | IDL             |
| ----------- | -------------------------| -------- | ---------------|
| api         | 对外 HTTP 服务接口        | cmd/api  | idl/api.thrift |
| note        | 内部的笔记微服务         | cmd/note | idl/note.thrift|
| user        | 内部的用户微服务         | cmd/user | idl/user.thrift|

在该项目中，note 服务和 user 服务先通过 ETCD 进行服务注册，API 服务再通过 ETCD 解析出它所依赖的服务的地址。微服务之间通过 RPC 进行通信，API 服务则通过 Hertz 框架对外提供 HTTP 接口。这种架构使得服务间通信更加高效可靠，同时也提高了系统的可扩展性和可维护性。note与user服务都通过gorm操作Mysql数据库进行CRUD。
## 安装与运行
本项目采用了hz与kitex进行代码生成
``` shell
go install github.com/cloudwego/hertz/cmd/hz@latest
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@v0.13.0
```
通过Makefile进行代码生成，下表为命令对照
| Catalog              | Command                             |
| --------------------| ----------------------------------- |
| hertz_api_model     | make hertz_gen_model                |
| hertz_api_client    | make hertz_gen_client               |
| kitex_user_client   | make kitex_gen_user                 |
| kitex_note_client   | make kitex_gen_note                 |
| hertz_api_new       | cd cmd/api && make hertz_new_api     |
| hertz_api_update    | cd cmd/api && make hertz_update_api  |
| kitex_user_server   | cd cmd/user && make kitex_gen_server |
| kitex_note_server   | cd cmd/note && make kitex_gen_server |

启动依赖环境
```bash
docker-compose up
```
启动user服务
```bash
cd cmd/user
make run
```
启动note服务
```bash
cd cmd/note
make run
```
启动api服务
```bash
cd cmd/api
make run
```
参考 `api_request/api_service/api_service_test.go` 文件，构建客户端，通过HTTP接口进行访问。该测试文件包含IDL中定义的所有api接口服务。
``` thrift
service ApiService {
    CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post="/v1/user/register")

    CheckUserResponse CheckUser(1: CheckUserRequest req) (api.post="/v1/user/login")

    CreateNoteResponse CreateNote(1: CreateNoteRequest req) (api.post="/v1/note")

    QueryNoteResponse QueryNote(1: QueryNoteRequest req) (api.get="/v1/note/query")

    UpdateNoteResponse UpdateNote(1: UpdateNoteRequest req) (api.put="/v1/note/:note_id")

    DeleteNoteResponse DeleteNote(1: DeleteNoteRequest req) (api.delete="/v1/note/:note_id")
}
```
## 开发指南：以User服务为例
```sh
cd cmd/user
make kitex_gen_server
```
即可看到生成的handler.go文件，内含user微服务的所有接口

在 `cmd/user/service/` 目录下对每个接口进行具体实现。
```sh
service/
├── check_user.go           # 用户登录、鉴权
├── create_user.go          # 创建用户
└── mget_user.go            # 获得多个用户信息 

```
查看 `cmd/user/dal/db/` 目录，可以看见user服务通过Gorm，从`consts`包中获得连接DSN，进行数据库连接。并且启用了`gormlogrus`作为logger，`opentelemetry`作为数据库访问的链路追踪。

由于gorm定义的model与kitex生成的rpc model并不完全相同，因此在 `cmd/user/pack` 包对user这个model进行了封装
```go
func User(u *db.User) *user.User {
	if u == nil {
		return nil
	}

	return &user.User{UserId: int64(u.ID), Username: u.Username, Avatar: "test"}
}
```
`main.go`
```golang
// 通过ETCD进行服务注册
r, err := etcd.NewEtcdRegistry([]string{consts.ETCDAddress})

// 设置OPTL云原生链路监测
p := provider.NewOpenTelemetryProvider(
	provider.WithServiceName(consts.UserServiceName),
        // Exporter地址
	provider.WithExportEndpoint(consts.ExportEndpoint),
	provider.WithInsecure(),
)
```
## 链路追踪、观测
根据docker-compose文件，我们启用了四个服务
1. otel-collector：OpenTelemetry Collector
2. jaeger-all-in-one：Jaeger分布式跟踪系统
3. victoriametrics：时间序列数据库
4. grafana：可视化仪表盘
otel-collector依赖于jaeger-all-in-one，因为在容器中启用了Jaeger的OTLP gRPC接收器。victoriametrics将端口8428暴露给其他服务，以便它们可以发送数据到时间序列数据库。grafana将端口3000暴露给其他服务，以便它们可以访问可视化仪表盘

以下为optl的部分配置信息
```yaml
receivers:
  otlp:

exporters:
  prometheusremotewrite:
    endpoint: "http://victoriametrics:8428/api/v1/write"

  logging:

  jaeger:
    endpoint: jaeger-all-in-one:14250

service:
  pipelines:
    traces:
      receivers: [ otlp ]
      processors: [ batch ]
      exporters: [ logging, jaeger ]
    metrics:
      receivers: [ otlp ]
      processors: [ batch ]
      exporters: [ logging, prometheusremotewrite ]
```
1. Receivers：首先定义了一个 OTLP Receiver，用于接收应用程序发送的数据。
2. Exporters：定义了三个 Exporters：Prometheus Remote Write Exporter（用于将数据导出到 VictoriaMetrics），Logging Exporter 和 Jaeger Exporter。
3. Service: 在这个配置文件中，定义了两个 Pipelines：Traces Pipeline（用于处理跟踪数据）和 Metrics Pipeline（用于处理指标数据），并且使用了上面定义的 Receivers、Processors 和 Exporters。
### Jaeger Grafana
- jaeger：浏览器打开 `http://127.0.0.1:16686/` 
![Jaeger](pic/jaeger.png)
- grafana：浏览器打开 `http://127.0.0.1:3000/`

登录 Grafana，并在左侧导航栏中选择 "Configuration"。

在 Configuration 页面上，选择 "Data Sources"，然后点击 "Add data source" 按钮。

在 Add data source 页面上，搜索并选择 Jaeger 数据源插件。

在 Jaeger 数据源的配置页面上，填写以下信息：

- Name：数据源名称
- URL：Jaeger 的地址 (由于我使用的wsl，需要设置为eth0的IP，而不是localhost)

配置完成后，点击 "Save & Test" 按钮来测试连接是否成功。如果成功，会出现 "Data source is working" 的提示。

配置完成后，在 Grafana 中创建一个新的 Dashboard，并选择 "Add Query" 按钮。

在 "Query" 页面中，选择 Jaeger 数据源，并选择要查询的 trace 数据（例如，按服务、操作、tag 等）。

完成 Query 配置后，就可以在 Dashboard 中看到 Jaeger 的 trace 数据。
![Grafana](pic/grafana.png)