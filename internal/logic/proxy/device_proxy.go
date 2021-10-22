package proxy

import (
	"context"
	"gim/pkg/pb"
)

type deviceProxy interface {
	ListOnlineByUserId(ctx context.Context, userId int64) ([]*pb.Device, error)
}

var DeviceProxy deviceProxy
