package endpoint

import (
	"errors"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

var (
	ErrInvalidRequest  = errors.New("request has only one class: Request")
	ErrResponseNull    = errors.New("response is null")
	ErrInvalidResponse = errors.New("response is invalid")
	ErrUnknownResponse = errors.New("unknown response")
)

type Request struct {
	RequestType string `json:"request_type"`
	Token       string `json:"token"`
	ServerName  string `json:"server_name"`
	Param       string `json:"param"`
	Data        []byte `json:"data"`
}

type Response struct {
	Result []byte `json:"result"`
	Error  error  `json:"error"`
}

type BifrostEndpoints struct {
	BifrostEndpoint     endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func NewBifrostEndpoints(svc service.Service) BifrostEndpoints {
	return BifrostEndpoints{
		BifrostEndpoint:     MakeBifrostServiceEndpoint(svc),
		HealthCheckEndpoint: MakeHealthCheckEndpoint(svc),
	}
}

func MakeBifrostServiceEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(Request); ok {
			r := service.NewRequest(ctx, req.RequestType, req.Token, req.ServerName, req.Param, req.Data)
			resp, err := svc.Deal(r)
			if err != nil {
				return Response{
					Result: nil,
					Error:  err,
				}, nil
			}
			if ret, err := resp.Bytes(); ret != nil {
				return Response{
					Result: ret,
					Error:  err,
				}, nil
			} else if watcher, err := resp.GetWatcher(); watcher != nil {
				return Watcher{
					dataChan:  watcher.GetDataChan(),
					errChan:   watcher.GetErrChan(),
					closeFunc: watcher.Close,
				}, err
			}
			return nil, ErrUnknownResponse
		}
		return nil, ErrInvalidRequest
	}
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthResponse{true}, nil
	}
}
