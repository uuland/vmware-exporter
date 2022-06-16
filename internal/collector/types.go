package collector

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
)

type Context struct {
	Context context.Context
	Client  *govmomi.Client
	Logger  *log.Logger
}

type Collector interface {
	Startup() error
	Shutdown() error
	Collectors() []prometheus.Collector
	Scrape(ctx *Context) error
}
