package fake

import (
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
	log "github.com/yongPhone/bifrost/pkg/log/v1"
)

type webServerStatus struct{}

func (w webServerStatus) Get(null *pbv1.Null, stream pbv1.WebServerStatus_GetServer) error {
	log.Infof("get web server status")

	return nil
}
