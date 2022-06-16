package collector

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
)

func NewService(cli *govmomi.Client, log *log.Logger) *Service {
	return &Service{cli: cli, log: log}
}

type Service struct {
	cs  []Collector
	cli *govmomi.Client
	log *log.Logger
}

func (s *Service) Start(feats string) error {
	cs, err := Registry.Load(strings.Split(feats, ",")...)
	if err != nil {
		return err
	}

	for _, c := range cs {
		cs2 := c.Collectors()
		for _, c2 := range cs2 {
			if err := prometheus.Register(c2); err != nil {
				return err
			}
		}

		if err := c.Startup(); err != nil {
			return err
		}
	}

	s.cs = cs
	return nil
}

func (s *Service) Stop() error {
	for _, c := range s.cs {
		if err := c.Shutdown(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Scrape() error {
	ctx1, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	ctx2 := &Context{
		Context: ctx1,
		Client:  s.cli,
		Logger:  s.log,
	}

	for _, c := range s.cs {
		if err := c.Scrape(ctx2); err != nil {
			return err
		}
	}

	return nil
}
