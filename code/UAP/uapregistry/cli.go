package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"uapregistry/config"
	"uapregistry/healthcheckmanager"
	rest "uapregistry/http"
	"uapregistry/leadermanager"
	"uapregistry/logger"
	"uapregistry/servicemanager"
	agent "uapregistry/storage/consulagent"
	"uapregistry/storage/consulagent/cache"
	"uapregistry/types"
	"uapregistry/utils"
)

// Name serviceName
const Name = "uapregistry"
const (
	// ExitCodeOK exit without error
	ExitCodeOK int = 0

	// ExitCodeError exit with error
	ExitCodeError = 10 + iota
	// ExitCodeInterrupt exit because of interrupt
	ExitCodeInterrupt
	// ExitCodeParseFlagsError exit because of parse flag error
	ExitCodeParseFlagsError
	// ExitCodeAgentError exit because of agent error
	ExitCodeAgentError
	// ExitCodeHTTPServerError exit because of server error
	ExitCodeHTTPServerError
	// ExitCodeRunnerError exit because of runner error
	ExitCodeRunnerError
	// ExitCodeConfigError exit because of config error
	ExitCodeConfigError
	ExitCodeWebHookError
)

// SignalLookup signal table
var SignalLookup = map[string]os.Signal{
	"SIGABRT": syscall.SIGABRT,
	"SIGALRM": syscall.SIGALRM,
	"SIGBUS":  syscall.SIGBUS,
	"SIGCHLD": syscall.SIGCHLD,
	//	"SIGCONT": syscall.SIGCONT,
	"SIGFPE": syscall.SIGFPE,
	"SIGHUP": syscall.SIGHUP,
	"SIGILL": syscall.SIGILL,
	"SIGINT": syscall.SIGINT,
	//	"SIGIO":   syscall.SIGIO,
	//	"SIGIOT":  syscall.SIGIOT,
	"SIGKILL": syscall.SIGKILL,
	"SIGPIPE": syscall.SIGPIPE,
	//	"SIGPROF": syscall.SIGPROF,
	"SIGQUIT": syscall.SIGQUIT,
	"SIGSEGV": syscall.SIGSEGV,
	//	"SIGSTOP": syscall.SIGSTOP,
	//	"SIGSYS":  syscall.SIGSYS,
	"SIGTERM": syscall.SIGTERM,
	"SIGTRAP": syscall.SIGTRAP,
	//	"SIGTSTP": syscall.SIGTSTP,
	//	"SIGTTIN": syscall.SIGTTIN,
	//	"SIGTTOU": syscall.SIGTTOU,
	"SIGURG": syscall.SIGURG,
	//  "SIGUSR1": syscall.SIGUSR1,
	//	"SIGUSR2": syscall.SIGUSR2,
	//	"SIGXCPU": syscall.SIGXCPU,
	//	"SIGXFSZ": syscall.SIGXFSZ,
}

// CLI - client
type CLI struct {
	sync.Mutex
	log                  logger.Logger
	outStream, errStream io.Writer
	servers              []*http.Server
	signalCh             chan os.Signal
	stopCh               chan struct{}
	stopped              bool
}

// NewCLI - create a new client
func NewCLI(out, err io.Writer) *CLI {
	return &CLI{
		log:       logger.GetLogger(),
		outStream: out,
		errStream: err,
		signalCh:  make(chan os.Signal, 1),
		stopCh:    make(chan struct{}),
	}
}

func (cli *CLI) startRoutines() {
	// health check
	if config.GetHealthCheckEnable() == "true" {
		logger.GetLogger().Info("HEALTH_CHECK_ENABLE is true,register self and start health manager")
		go registerSelf()
		go healthcheckmanager.StartHealthCheckManager()
	}
}

func registerSelf() {
	serviceIP := utils.GetNodeIP()
	if serviceIP == "" {
		logger.GetLogger().Errorf("failed to register self:service ip not found")
		return
	}

	port, err := strconv.Atoi(config.GetHTTPListenPort())
	if err != nil {
		logger.GetLogger().Errorf("failed to register self:invalid http listen port %s", config.GetHTTPListenPort())
		return
	}

	svc := types.Service{
		Name:     "uapregistry-health-check",
		Protocol: "TCP",
		Host:     serviceIP,
		Port:     port,
		PersistentCheck: &types.PersistentCheckInfo{
			CheckType:    "http",
			CheckHTTPURL: "/health",
		},
	}
	utils.FillServiceDefaultValue(&svc)

	for {
		_, _, err := servicemanager.NewServiceManager().PostService(&svc, true)
		if err != nil {
			logger.GetLogger().Errorf("failed to register self:%v,retry after 5 seconds ", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}

// Run - run a cli
func (cli *CLI) Run(args []string) int {
	defer cli.shutdown()
	// Init Config
	config.InitConfig()
	// Parse the flags
	agentCfg, httpCfg, err := cli.ParseFlags(args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return cli.handleError(err, ExitCodeParseFlagsError)
	}
	cli.showConfigs(agentCfg, httpCfg)

	// Create Consul Agent and set watch package var
	err = agent.InitLocalAgent(agentCfg)
	if err != nil {
		return cli.handleError(err, ExitCodeAgentError)
	}

	go cache.InitServiceCache()
	go cache.InitRouteCache()

	if config.GetHealthCheckEnable() == "true" {
		leadermanager.StartLeaderManager(agent.GetLocalAgent().GetConsulClient(), utils.GetNodeIP())
	}

	cli.log.Info("start http server on %s", httpCfg.GetHTTPPort())
	// Launch the HTTP Server
	httpCh := make(chan error)
	cli.servers = append(cli.servers, rest.StartHTTPServer(httpCfg, httpCh)...)

	cli.startRoutines()

	// Listen for signals
	return cli.handleFin(httpCh)
}

func (cli *CLI) handleFin(ch chan error) int {
	signal.Notify(cli.signalCh)
	for {
		select {
		case e := <-ch:
			cli.log.Errorf("HTTP Server exited with error:%v", e)
			buf := make([]byte, 1<<20)
			stackSize := runtime.Stack(buf, true)
			cli.log.Infof("%s\n", string(buf[0:stackSize]))
			cli.log.Flush()
			return cli.handleError(e, ExitCodeHTTPServerError)
		case s := <-cli.signalCh:
			if code := cli.handleSignals(s); code > 0 {
				return code
			}
		case <-cli.stopCh:
			return ExitCodeOK
		}
	}
}

func (cli *CLI) handleSignals(s os.Signal) int {
	if s != SignalLookup["SIGURG"] {
		cli.log.Warnf("cli receiving signal %q", s)
	}

	switch s {
	case SignalLookup["SIGINT"]:
		fallthrough
	case SignalLookup["SIGTERM"]:
		fallthrough
	case SignalLookup["SIGKILL"]:
		cli.log.Flush()
		return ExitCodeInterrupt
	case SignalLookup["SIGQUIT"]:
		buf := make([]byte, 1<<20)
		stackSize := runtime.Stack(buf, true)
		cli.log.Infof("%s\n", string(buf[0:stackSize]))
	case SignalLookup["SIGCHLD"]:
		cli.reapChild()
	default:
	}
	return 0
}

func (cli *CLI) reapChild() {
	for {
		var status syscall.WaitStatus
		cpid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
		if err != nil {
			if err != syscall.ECHILD {
				cli.log.Warnf("wait4 after SIGCHLD: %v", err)
			}
			break
		}
		if cpid < 1 {
			break
		}
		if status.Exited() {
			cli.log.Warnf("Reaped process with pid %d, exited with status %d", cpid, status.ExitStatus())
		} else if status.Signaled() {
			cli.log.Warnf("Reaped process with pid %d, exited on %s", cpid, status.Signal())
		} else {
			cli.log.Warnf("Reaped process with pid %d", cpid)
		}
	}
}

// stop - stop a client
func (cli *CLI) stop() {
	cli.Lock()
	defer cli.Unlock()

	if cli.stopped {
		return
	}

	close(cli.stopCh)
	cli.stopped = true
}

func (cli *CLI) shutdown() {
	if len(cli.servers) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, value := range cli.servers {
		err := value.Shutdown(ctx)
		if err != nil {
			cli.log.Errorf("Failed to shutdown HTTP Server:%v", err)
		}
	}
}

// ParseFlags - parse input params
func (cli *CLI) ParseFlags(args []string) (*agent.Config, *rest.Config, error) {

	agentCfg := agent.DefaultConfig()
	httpCfg := rest.DefaultConfig()

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() { fmt.Fprintf(cli.errStream, usage, Name) }

	flags.Var((funcVar)(func(s string) error {
		agentCfg.SetConsulAgent(s)
		return nil
	}), "consul-agent", "")

	flags.Var((funcVar)(func(s string) error {
		httpCfg.SetHTTPPort(s)
		return nil
	}), "listen-port", "")

	flags.Var((funcVar)(func(s string) error {
		httpCfg.SetHTTPBindIP(s)
		return nil
	}), "bind-ip", "")

	// If there was a parser error, stop
	if err := flags.Parse(args); err != nil {
		return nil, nil, err
	}

	return agentCfg, httpCfg, nil
}

func (cli *CLI) handleError(err error, status int) int {
	fmt.Fprintf(cli.errStream, "||%s|| Service Discovery Client returned errors:\n%v\n", time.Now().String(), err)
	return status
}

func (cli *CLI) showConfigs(acfg *agent.Config, hcfg *rest.Config) {
	cli.log.Warn("Starting uapregistry:  ...")

	cli.log.Warn("ConsulAgent Addr:", acfg.GetConsulAgent())
	cli.log.Warnf("Bind IP:%v", hcfg.GetHTTPIPs())
	cli.log.Warn("Listen Port:", hcfg.GetHTTPPort())
	cli.log.Warn("HEALTH_CHECK_ENABLE:", config.GetHealthCheckEnable())
	cli.log.Flush()
}

const usage = `
Usage: %s [options]

  A Service Discovery Client tool using Consul as the backend.

Options:

  -consul-agent=<address>
      Sets the address of the Consul instance
  -listen-port=<port>
      Sets the listen port of the http server
`
