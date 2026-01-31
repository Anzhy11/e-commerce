package userService

import (
	"fmt"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/rs/zerolog"
)

type UserServiceInterface interface {
	GetUser(id string) string
}

type userService struct {
	cfg *config.Config
	log *zerolog.Logger
}

func New(cfg *config.Config, log *zerolog.Logger) UserServiceInterface {
	return &userService{
		cfg: cfg,
		log: log,
	}
}

func (s *userService) GetUser(id string) string {
	s.log.Info().Msg(fmt.Sprintf("Hit user service %s", s.cfg.Server.Port))
	return "user service here"
}
