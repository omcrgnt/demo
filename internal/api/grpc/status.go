package grpc

import (
	"errors"

	"github.com/omcrgnt/demo/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ServiceError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case err == nil:
		return nil
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
