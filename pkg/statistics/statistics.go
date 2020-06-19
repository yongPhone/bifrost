package statistics

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"strconv"
)

func HTTPServersNum(ctx resolv.Context) int {
	http := resolv.GetHTTP(ctx)
	if http == nil {
		return 0
	}
	return len(http.Servers())
}

func HTTPPorts(ctx resolv.Context) []int {
	var ports []int
	http := resolv.GetHTTP(ctx)
	if http == nil {
		return nil
	}
	for _, parser := range resolv.GetPorts(http) {
		port, err := strconv.Atoi(parser.(*resolv.Key).Value)
		if err != nil {
			continue
		}
		ports = resolv.SortInsertUniqInt(ports, port)
	}
	return ports
}

func HTTPServerNames(ctx resolv.Context) (serverNames []string) {
	for _, parser := range resolv.GetHTTPServers(ctx, resolv.ServerName) {
		if serverNameKey := resolv.GetServerName(parser.(*resolv.Server)); serverNameKey != nil {
			serverNames = resolv.AppendNewString(serverNames, serverNameKey.(*resolv.Key).Value)
		}
	}
	return
}

func HTTPLocationsNum(ctx resolv.Context) int {
	return len(resolv.GetLocations(ctx))
}

func StreamServersNum(ctx resolv.Context) int {
	stream := resolv.GetStream(ctx)
	if stream == nil {
		return 0
	}
	return len(stream.Servers())
}

func StreamPorts(ctx resolv.Context) []int {
	//var ports []string
	var ports []int
	stream := resolv.GetStream(ctx)
	if stream == nil {
		return nil
	}
	for _, parser := range resolv.GetPorts(stream) {
		//ports = appendNewString(ports, parser.(*resolv.Key).Value)
		port, err := strconv.Atoi(parser.(*resolv.Key).Value)
		if err != nil {
			continue
		}
		ports = resolv.SortInsertUniqInt(ports, port)
	}
	return ports
}
