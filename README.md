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
参考 `api_request/api_service/api_service_test.go` 文件，构建客户端，通过HTTP接口进行访问。
## 开发指南