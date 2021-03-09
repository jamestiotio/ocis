package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/thejerf/suture"

	onlyoffice "github.com/owncloud/ocis/onlyoffice/pkg/command"
	settings "github.com/owncloud/ocis/settings/pkg/command"

	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/owncloud/ocis/ocis/pkg/runtime/log"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	halt = make(chan os.Signal, 1)
	done = make(chan struct{}, 1)
)

// Service represents a RPC service.
type Service struct {
	Supervisor       *suture.Supervisor
	ServicesRegistry map[string]func(context.Context, *ociscfg.Config) suture.Service
	serviceToken     map[string][]suture.ServiceToken
	context          context.Context
	cancel           context.CancelFunc
	Log              zerolog.Logger
	wg               *sync.WaitGroup
	done             bool
	cfg              *ociscfg.Config
}

// TODO(refs) think of a less confusing naming
type RuntimeSutureService struct {
	Context    context.Context
	CancelFunc context.CancelFunc
}

// loadFromEnv would set cmd global variables. This is a workaround spf13/viper since pman used as a library does not
// parse flags.
func loadFromEnv() (*config.Config, error) {
	cfg := config.NewConfig()
	viper.AutomaticEnv()
	if err := viper.BindEnv("port", "RUNTIME_PORT"); err != nil {
		return nil, err
	}
	if viper.GetString("port") != "" {
		cfg.Port = viper.GetString("port")
	}

	return cfg, nil
}

// NewService returns a configured service with a controller and a default logger.
// When used as a library, flags are not parsed, and in order to avoid introducing a global state with init functions
// calls are done explicitly to loadFromEnv().
// Since this is the public constructor, options need to be added, at the moment only logging options
// are supported in order to match the running OwnCloud services structured log.
func NewService(options ...Option) (*Service, error) {
	opts := NewOptions()

	for _, f := range options {
		f(opts)
	}

	l := log.NewLogger(
		log.WithPretty(opts.Log.Pretty),
	)

	globalCtx, cancelGlobal := context.WithCancel(context.Background())

	s := &Service{
		ServicesRegistry: make(map[string]func(context.Context, *ociscfg.Config) suture.Service),
		Log:              l,

		serviceToken: make(map[string][]suture.ServiceToken),
		context:      globalCtx,
		cancel:       cancelGlobal,
		wg:           &sync.WaitGroup{},
		cfg:          opts.Config,
	}

	s.ServicesRegistry["onlyoffice"] = onlyoffice.NewSutureService
	s.ServicesRegistry["settings"] = settings.NewSutureService

	return s, nil
}

// Start an rpc service.
func Start(o ...Option) error {
	s, err := NewService(o...)
	if err != nil {
		if s != nil {
			s.Log.Fatal().Err(err)
		}
	}

	s.Supervisor = suture.NewSimple("ocis")

	if err := rpc.Register(s); err != nil {
		if s != nil {
			s.Log.Fatal().Err(err)
		}
	}
	rpc.HandleHTTP()

	signal.Notify(halt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// TODO(refs) change default port
	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "localhost", "6060"))
	if err != nil {
		s.Log.Fatal().Err(err)
	}

	// handle panic within the Service scope.
	defer func() {
		if r := recover(); r != nil {
			reason := strings.Builder{}
			// TODO(refs) change default port
			if _, err := net.Dial("localhost", "6060"); err != nil {
				reason.WriteString("runtime address already in use")
			}

			fmt.Println(reason.String())
		}
	}()
	for k, _ := range s.ServicesRegistry {
		s.serviceToken[k] = append(s.serviceToken[k], s.Supervisor.Add(s.ServicesRegistry[k](s.context, s.cfg)))
	}

	go s.Supervisor.ServeBackground()
	go trap(s)

	return http.Serve(l, nil)
}

// Start indicates the Service Controller to start a new supervised service as an OS thread.
func (s *Service) Start(name string, reply *int) error {
	if _, ok := s.ServicesRegistry[name]; !ok {
		*reply = 1
		return nil
	}
	//s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(s.ServicesRegistry[name]))
	//s.serviceToken["settings"] = append(s.serviceToken[name], s.Supervisor.Add(settings.NewSutureService(s.context, s.cfg)))
	*reply = 0
	return nil
}

// List running processes for the Service Controller.
func (s *Service) List(args struct{}, reply *string) error {
	return nil
}

// Kill a supervised process by subcommand name.
func (s *Service) Kill(name string, reply *int) error {
	if len(s.serviceToken[name]) > 0 {
		for i := range s.serviceToken[name] {
			fmt.Printf("\n\n%s\n%+v\n\n", name, s.serviceToken[name])
			if err := s.Supervisor.Remove(s.serviceToken[name][i]); err != nil {
				return err
			}
			delete(s.serviceToken, name)
		}
	} else {
		return fmt.Errorf("service %s not found", name)
	}

	return nil
}

// trap blocks on halt channel. When the runtime is interrupted it
// signals the controller to stop any supervised process.
func trap(s *Service) {
	<-halt
	s.done = true
	s.wg.Wait()
	s.cancel()
	s.Log.Debug().Str("service", "runtime service").Msgf("terminating with signal: %v", s)
	close(done)
	os.Exit(0)
}
