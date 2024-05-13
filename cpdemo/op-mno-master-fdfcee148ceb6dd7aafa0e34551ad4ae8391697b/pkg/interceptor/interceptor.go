package interceptor

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"time"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		// 从传入上下文获取元数据 暂不需要，先保留
		//md, ok := metadata.FromIncomingContext(ctx)
		//if !ok {
		//	return nil, fmt.Errorf("couldn't parse incoming context metadata")
		//}

		// 获取客户端IP地址
		ip, err := getClientIP(ctx)
		if err != nil {
			return nil, err
		}
		m, err := handler(ctx, req)
		end := time.Now()
		log.Printf("请求方法:【%s】，客户端ip:【%v】，请求参数:【%v】,请求时间:【%s】，响应时间:【%s】，耗时:【%v】，err:【%v】", info.FullMethod, ip, req, utils.FormatToDateTimeStr(start), utils.FormatToDateTimeStr(end), end.Sub(start), err)
		return m, err
	}
}

type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	fmt.Printf("Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	fmt.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		ctx := ss.Context()
		// 获取客户端IP地址
		ip, err := getClientIP(ctx)
		if err != nil {
			return err
		}
		// RPC 方法真正执行的逻辑
		// 调用RPC方法(invoking RPC method)
		err = handler(srv, newWrappedStream(ss))
		end := time.Now()
		log.Printf("请求方法:【%s】，客户端ip:【%v】，请求时间:【%s】，响应时间:【%s】，耗时:【%v】，err:【%v】", info.FullMethod, ip, utils.FormatToDateTimeStr(start), utils.FormatToDateTimeStr(end), end.Sub(start), err)
		return err
	}
}

func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("couldn't parse client IP address")
	}
	return p.Addr.String(), nil
}
