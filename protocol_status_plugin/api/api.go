package api

import (
	"context"
	pluginpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/plugin"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	manager_presenter "vatz-plugin-near/protocol_status/manager"
	policy "vatz-plugin-near/protocol_status/policy"
)

var (
	ManagerInstance manager_presenter.Manager
	policyExecutor  policy.Executor
)

type GrpcService struct {
	pluginpb.UnimplementedManagerPluginServer
}

func (s *GrpcService) Init(context.Context, *emptypb.Empty) (*pluginpb.PluginInfo, error) {
	// TODO: TBD
	return nil, nil
}

func (s *GrpcService) Verify(context.Context, *emptypb.Empty) (*pluginpb.VerifyInfo, error) {
	// TODO: TBD
	return nil, nil
}

func (s *GrpcService) Execute(ctx context.Context, req *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {

	log.Println("pluginServer.Execute")
	resp, err := manager_presenter.RunManager().Execute(ctx, req)

	if err != nil {
		return &pluginpb.ExecuteResponse{State: pluginpb.ExecuteResponse_FAILURE}, nil
	}

	return resp, nil
}
