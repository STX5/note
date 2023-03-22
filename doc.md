1. thrift版本需要设置为v0.13.0，hz自动生成的代码不兼容最新版thrift
2. 需要安装验证器: Validator 是用于支持结构体校验能力的 thriftgo 插件 `go install github.com/cloudwego/thrift-gen-validator@latest`
3. hz model/ hz client/ hz new 有区别
4. hz gen的文件编译：
    ```
    $ sh build.sh
    ```
    执行上述命令后，会生成一个 output 目录，里面含有编译产物。

    运行：
    ```
    $ sh output/bootstrap.sh
    ```
5. 创建cmd/api目录，创建Makefile。make hertz_new_api，
    然后生成hertz_handler/api下的api_service.go文件即为业务逻辑的handler，包含所有定义的方法。
    修改hertz_router的middleware.go
6. JWT Middleware 因为 JWT 的核心是认证与授权，所以在使用 Hertz 的 
   jwt 扩展时，
   不仅需要为 /login 接口绑定认证逻辑 authMiddleware.LoginHandler。还要以中间件的方式，
   为需要授权访问的路由组注入授权逻辑 authMiddleware.MiddlewareFunc()