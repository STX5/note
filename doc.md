1. thrift版本需要设置为v0.13.0，hz自动生成的代码不兼容最新版thrift
2. 需要安装验证器: Validator 是用于支持结构体校验能力的 thriftgo 插件 `go install github.com/cloudwego/thrift-gen-validator@latest`
3. hz gen的文件编译：
    ```
    $ sh build.sh
    ```
    执行上述命令后，会生成一个 output 目录，里面含有编译产物。

    运行：
    ```
    $ sh output/bootstrap.sh
    ```
4. 创建cmd/api目录，创建Makefile。make hertz_new_api，
    然后生成hertz_handler/api下的api_service.go文件即为业务逻辑的handler，包含所有定义的方法。