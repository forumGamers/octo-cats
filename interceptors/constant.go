package interceptors

import (
	"context"

	"github.com/forumGamers/octo-cats/pkg/user"
	"google.golang.org/grpc"
)

type Interceptor interface {
	UnaryAuthentication(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	GetUserFromCtx(ctx context.Context) user.User
	Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
}

type InterceptorImpl struct{}

func NewInterCeptor() Interceptor {
	return &InterceptorImpl{}
}
