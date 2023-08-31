package gapi

import (
	"context"
	zlog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net/http"
	"time"
)

func GrpcLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {

	now := time.Now()
	rsp, err := handler(ctx, req)
	duration := time.Since(now)

	statusCode := codes.Unknown
	/*
		if st, ok := statusCode.FromErr(err); ok {
			statusCode = st.Code()
		}
	*/

	logger := zlog.Info()
	if logger != nil {
		logger = zlog.Error().Err(err)
	}

	logger.
		Str("protocol", "grpc").
		Int("status-code", int(statusCode)).
		Str("status-text", statusCode.String()).
		Str("method", info.FullMethod).
		Dur("Duration", duration).
		Msg("received a gRPC request")

	return rsp, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: rsp,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rsp, req)
		duration := time.Since(startTime)

		logger := zlog.Info()
		if rec.StatusCode != http.StatusOK {
			logger = zlog.Error().Bytes("body", rec.Body)
		}

		logger.
			Str("protocol", "http").
			Str("status-code", req.Method).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("Duration", duration).
			Msg("received a http request")
	})
}
