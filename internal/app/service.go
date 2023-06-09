package app

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"ozon-test-unzhakov/internal/dto"
	"ozon-test-unzhakov/internal/service"
	desc "ozon-test-unzhakov/pkg"
)

type MicroserviceServer struct {
	desc.UnimplementedMicroserviceServer
	linkService service.LinkService
}

func NewMicroservice(ls service.LinkService) *MicroserviceServer {
	return &MicroserviceServer{linkService: ls}
}

func (m *MicroserviceServer) CreateShortLink(ctx context.Context, r *desc.CreateShortLinkRequest) (*desc.CreateShortLinkResponse, error) {
	shorted, err := m.linkService.CreateShortLink(&dto.Link{Link: r.Link})
	if err != nil {
		return nil, err
	}
	return &desc.CreateShortLinkResponse{Code: shorted.Code}, nil
}

func (m *MicroserviceServer) GetInitialLink(ctx context.Context, r *desc.GetInitialLinkRequest) (*desc.GetInitialLinkResponse, error) {
	initial, err := m.linkService.GetInitialLink(&dto.Link{Code: r.Code})
	if err != nil {
		return nil, err
	}
	return &desc.GetInitialLinkResponse{Link: initial.Link}, nil
}

func (m *MicroserviceServer) RedirectToInitialLink(ctx context.Context, r *desc.RedirectToInitialLinkRequest) (*desc.EmptyMessage, error) {
	initial, err := m.linkService.GetInitialLink(&dto.Link{Code: r.Code})
	if err != nil {
		return nil, err
	}
	header := metadata.Pairs("Location", initial.Link)
	err = grpc.SendHeader(ctx, header)
	if err != nil {
		return nil, err
	}
	return &desc.EmptyMessage{}, nil
}
