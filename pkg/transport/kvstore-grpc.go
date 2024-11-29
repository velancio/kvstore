package transport

import (
	"censys/internal/kvstore"
	"censys/proto/gen/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KvStoreServer is a struct that implements the KvStoreServiceServer interface
type KvStoreServer struct {
	proto.UnimplementedKvStoreServiceServer
	Store kvstore.KeyValueStore
}

// Get returns the value for the given key
func (s *KvStoreServer) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetResponse, error) {
	value, ok := s.Store.Get(ctx, request.GetKey())
	if !ok {
		return &proto.GetResponse{
			Value:   value,
			Success: false,
		}, status.Errorf(codes.NotFound, "key not found")
	}

	return &proto.GetResponse{
		Value:   value,
		Success: true,
	}, nil
}

// Set sets the value for the given key
func (s *KvStoreServer) Set(ctx context.Context, request *proto.SetRequest) (*proto.SetResponse, error) {
	err := s.Store.Set(ctx, request.GetKey(), request.GetValue())
	if err != nil {
		if err.Error() == "key cannot be empty" {
			return &proto.SetResponse{
				Success: false,
			}, status.Errorf(codes.InvalidArgument, "key cannot be empty")
		}
		return &proto.SetResponse{
			Success: false,
		}, err
	}

	return &proto.SetResponse{
		Success: true,
	}, nil
}

// Delete deletes the value for the given key
func (s *KvStoreServer) Delete(ctx context.Context, request *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	keyToDelete := request.GetKey()
	_, ok := s.Store.Get(ctx, keyToDelete)
	if !ok {
		return &proto.DeleteResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "key not found")
	}

	err := s.Store.Delete(ctx, keyToDelete)
	if err != nil {
		return &proto.DeleteResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "failed to delete key not found %#v", err)
	}

	return &proto.DeleteResponse{
		Success: true,
	}, nil
}
