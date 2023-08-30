package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGateWayUserAgentHeader = "grpcgateway-user-agent"
	grpcGateWayClientIPHeader  = "x-forwareded-for"
	userAgentHeader            = "user-agent"
)

type MetaData struct {
	ClientIP  string
	UserAgent string
}

func (s *Server) extractMetaData(ctx context.Context) *MetaData {
	mtd := &MetaData{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if userAgent := md.Get(grpcGateWayUserAgentHeader); len(userAgent) > 0 {
			mtd.UserAgent = userAgent[0]
		}
		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			mtd.ClientIP = userAgent[0]
		}
		if clientIP := md.Get(grpcGateWayClientIPHeader); len(clientIP) > 0 {
			mtd.ClientIP = clientIP[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtd.ClientIP = p.Addr.String()
	}

	return mtd
}
