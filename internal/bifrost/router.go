package bifrost

import (
	"time"

	epv1 "github.com/yongPhone/bifrost/internal/bifrost/endpoint/v1"
	"github.com/yongPhone/bifrost/internal/bifrost/middleware"
	svcv1 "github.com/yongPhone/bifrost/internal/bifrost/service/v1"
	storev1 "github.com/yongPhone/bifrost/internal/bifrost/store/v1"
	txpv1 "github.com/yongPhone/bifrost/internal/bifrost/transport/v1"
	handlerv1 "github.com/yongPhone/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/yongPhone/bifrost/internal/bifrost/transport/v1/options"
	genericgrpcserver "github.com/yongPhone/bifrost/internal/pkg/server"
	log "github.com/yongPhone/bifrost/pkg/log/v1"
)

func initRouter(server *genericgrpcserver.GenericGRPCServer) {
	svcIns := initService()
	initMiddleware(&svcIns)
	initController(svcIns, server)
}

func initService() svcv1.ServiceFactory {
	storeIns := storev1.Client()

	return svcv1.NewServiceFactory(storeIns)
}

func initMiddleware(svc *svcv1.ServiceFactory) {
	middlewaresIns := middleware.GetMiddlewares()
	for name, m := range middlewaresIns {
		log.Infof("Install middleware: %s", name)
		*svc = m(*svc)
	}
}

func initController(svc svcv1.ServiceFactory, server *genericgrpcserver.GenericGRPCServer) {
	// v1 transport
	eps := epv1.NewEndpoints(svc)
	hs := handlerv1.NewHandlersFactory(eps)
	opts := &options.Options{
		ChunkSize:          server.ChunkSize,
		RecvTimeoutMinutes: server.ReceiveTimeout / time.Minute,
	}

	txp := txpv1.New(hs, opts)
	{
		// register bifrost services
		registers := txpv1.NewBifrostServiceRegister(txp)
		server.RegisterServices(registers.Generate())
	}
}
