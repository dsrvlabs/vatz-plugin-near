package manager

import (
	"context"
	pluginpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/plugin"
	"google.golang.org/protobuf/types/known/emptypb"
	worker_presenter "vatz-plugin-near/protocol_status/worker"
)

type Manager interface {
	Init(context.Context, *emptypb.Empty) (*pluginpb.PluginInfo, error)
	Verify(context.Context, *emptypb.Empty) (*pluginpb.VerifyInfo, error)
	Execute(ctx context.Context, req *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error)
}

func RunManager() *worker_presenter.ManagerWorker {
	return &worker_presenter.ManagerWorker{}
}
