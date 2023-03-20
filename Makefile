# http interface
hertz_gen_model:
	hz model --idl=idl/api.thrift --mod=note --model_dir=hertz_gen
hertz_gen_client:
	hz client --idl=idl/api.thrift --base_domain=http://127.0.0.1:8080 --client_dir=api_request --mod=note --model_dir=hertz_gen
# rpc interface
kitex_gen_user:
	kitex --thrift-plugin validator -module note idl/user.thrift # execute in the project root directory
kitex_gen_note:
	kitex --thrift-plugin validator -module note idl/note.thrift # execute in the project root directory

# not compatible with thrift@latest(0.18.0) 
install_hz_latest:
	go install github.com/cloudwego/hertz/cmd/hz@latest
