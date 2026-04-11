package script

import (
	"github.com/dop251/goja"
	"github.com/ecumeurs/upsiloncli/internal/api"
	"github.com/ecumeurs/upsiloncli/internal/display"
	"github.com/ecumeurs/upsiloncli/internal/endpoint"
	"github.com/ecumeurs/upsiloncli/internal/session"
	"github.com/ecumeurs/upsiloncli/internal/ws"
	"io"
)

type Agent struct {
	ID           string
	Session      *session.Session
	Client       *api.Client
	Listener     *ws.Listener
	Registry     *endpoint.Registry
	VM           *goja.Runtime
	Logger       io.Writer
	TeardownHook goja.Callable
	Shared       *SharedStore
}

func NewAgent(id, baseURL string, reg *endpoint.Registry, logger io.Writer, shared *SharedStore) *Agent {
	sess := session.New()
	printer := display.NewPrinterWithWriter(logger)
	client := api.NewClient(baseURL, sess, printer)
	
	agent := &Agent{
		ID:       id,
		Session:  sess,
		Client:   client,
		Listener: ws.NewListener(client, sess, printer),
		Registry: reg,
		VM:       goja.New(),
		Logger:   logger,
		Shared:   shared,
	}

	agent.bindJSAPI()
	return agent
}
