package app

import (
	"context"
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

func (m *MicroserviceServer) ShortLink(ctx context.Context, r *desc.CreateShortLinkRequest) (*desc.CreateShortLinkResponse, error) {
	shorted, err := m.linkService.CreateShortLink(r.Link)
	if err != nil {
		return nil, err
	}
	return &desc.CreateShortLinkResponse{Code: shorted.Link}, nil
}

func (m *MicroserviceServer) FullLink(ctx context.Context, r *desc.GetInitialLinkRequest) (*desc.GetInitialLinkResponse, error) {
	initial, err := m.linkService.GetInitialLink(r.Code)
	if err != nil {
		return nil, err
	}
	return &desc.GetInitialLinkResponse{Link: initial.Link}, nil
}
