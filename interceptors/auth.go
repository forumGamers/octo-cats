package interceptors

import (
	"context"
	"os"

	"github.com/forumGamers/octo-cats/pkg/user"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (i *InterceptorImpl) UnaryAuthentication(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := metadata["access_token"]
	if len(values) < 1 {
		return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
	}

	claim := jwt.MapClaims{}
	if token, err := jwt.ParseWithClaims(values[0], &claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}); err != nil || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
	}

	return handler(context.WithValue(ctx, "user", claim), req)
}

func (i *InterceptorImpl) GetUserFromCtx(ctx context.Context) user.User {
	var user user.User

	claim, ok := ctx.Value("user").(jwt.MapClaims)
	if !ok {
		return user
	}

	for key, val := range claim {
		switch key {
		case "UUID":
			user.UUID = val.(string)
		case "loggedAs":
			user.LoggedAs = val.(string)
		default:
			break
		}
	}
	return user
}
