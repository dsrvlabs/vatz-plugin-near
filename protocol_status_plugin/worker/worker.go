package worker

import (
	"context"
	pluginpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/plugin"
	managerpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"os/exec"
	"strings"
)

type ManagerWorker struct{}

func (p *ManagerWorker) Init(context.Context, *emptypb.Empty) (*pluginpb.PluginInfo, error) {
	return nil, nil
}

func (p *ManagerWorker) Verify(ctx context.Context, in *managerpb.VerifyRequest) (*managerpb.VerifyInfo, error) {
	return nil, nil
}

func (p *ManagerWorker) Execute(ctx context.Context, req *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {
	//TODO
	// 이곳을 기점으로 어떻게 플러그인 기준을 어떻게 MVC 기준으로 가져갈지 생각할 수 있게 한다.?
	// 아니면 그것은 자율에 맞긴다.

	c, b := exec.Command("ps -ef | grep near"), new(strings.Builder)
	c.Stdout = b
	c.Run()

	return &pluginpb.ExecuteResponse{
		State:        pluginpb.ExecuteResponse_SUCCESS,
		Message:      b.String(),
		ResourceType: "NearProtocol",
	}, nil
}
