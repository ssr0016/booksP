package registry

import "context"

// type Service interface {
// 	Run(ctx context.Context) error
// }

type RunFunc func(ctx context.Context) error

type Service struct {
	Services []RunFunc
}

func NewServiceRegistry(services ...RunFunc) *Service {
	return &Service{services}
}

func (s *Service) GetServices() []RunFunc {
	return s.Services
}
