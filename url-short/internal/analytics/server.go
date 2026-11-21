
package analytics

import (
	"context"

	analyticsv1 "url-shortener/gen/go/analytics/v1"
	"url-shortener/internal/store/db"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	analyticsv1.UnimplementedAnalyticsServiceServer
	queries *db.Queries
}

func NewServer(queries *db.Queries) *Server {
	return &Server{
		queries: queries,
	}
}

func (s *Server) RecordClick(
	ctx context.Context,
	req *analyticsv1.RecordClickRequest,
) (*analyticsv1.RecordClickResponse, error) {

	if req.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "URL code is required")
	}

	_, err := s.queries.GetURLByCode(ctx, req.GetCode())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, "URL not found")
		}

		return nil, status.Error(codes.Internal, "failed to find URL")
	}

	err = s.queries.RecordClick(ctx, req.GetCode())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to record click")
	}

	return &analyticsv1.RecordClickResponse{
		Success: true,
	}, nil
}

func (s *Server) GetURLStats(
	ctx context.Context,
	req *analyticsv1.GetURLStatsRequest,
) (*analyticsv1.GetURLStatsResponse, error) {

	if req.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "URL code is required")
	}

	_, err := s.queries.GetURLByCode(ctx, req.GetCode())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, "URL not found")
		}

		return nil, status.Error(codes.Internal, "failed to find URL")
	}

	clicks, err := s.queries.GetURLClickCount(ctx, req.GetCode())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve click statistics")
	}

	return &analyticsv1.GetURLStatsResponse{
		Code:       req.GetCode(),
		TotalClicks: clicks,
	}, nil
}
