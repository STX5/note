package rpc

import (
	"context"
	"note/kitex_gen/user"
	"note/kitex_gen/user/userservice"
	"note/pkg/consts"
	"note/pkg/errno"

	cw "note/pkg/commonMiddleware"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"gorm.io/plugin/opentelemetry/provider"
)

var userClient userservice.Client

func initUser() {
	r, err := etcd.NewEtcdResolver([]string{consts.ETCDAddress})
	if err != nil {
		panic(err)
	}
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(consts.NoteServiceName),
		provider.WithExportEndpoint(consts.ExportEndpoint),
		provider.WithInsecure(),
	)
	defer func(ctx context.Context, p provider.OtelProvider) {
		_ = p.Shutdown(ctx)
	}(context.Background(), p)

	c, err := userservice.NewClient(
		consts.UserServiceName, // DestService
		client.WithResolver(r),
		client.WithMuxConnection(1),
		client.WithMiddleware(cw.CommonMiddleware),
		client.WithInstanceMW(cw.ClientMiddleware),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.NoteServiceName}),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

// MGetUser multiple get list of user info
func MGetUser(ctx context.Context, req *user.MGetUserRequest) (map[int64]*user.User, error) {
	resp, err := userClient.MGetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.BaseResp.StatusCode, resp.BaseResp.StatusMessage)
	}
	res := make(map[int64]*user.User)
	for _, u := range resp.Users {
		res[u.UserId] = u
	}
	return res, nil
}
