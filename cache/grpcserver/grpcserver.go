package grpcserver

import (
	"context"

	"github.com/tjroeder/censys_take_home/cache/internal/cache"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the gRPC CacheServer, and delegates
// storage to the internal CacheService
type Server struct {
	Cache cache.CacheService
	UnimplementedCacheServer
}

// New initializes a new gRPC Caching Server
func New(c cache.CacheService) *Server {
	return &Server{
		Cache: c,
	}
}

// Get retrieves cached record matching given key
func (s *Server) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	resp, ok := s.Cache.Get(req.Key)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cache: no key %s found", req.Key)
	}
	return &GetResponse{Value: resp}, nil
}

// Set unconditionally adds a key-value pair record into the cache
func (s *Server) Set(ctx context.Context, req *SetRequest) (*SetResponse, error) {
	s.Cache.Set(req.Key, req.Value)
	return &SetResponse{}, nil
}

// Delete removes cached record matching given key
func (s *Server) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	s.Cache.Delete(req.Key)
	return &DeleteResponse{}, nil
}
