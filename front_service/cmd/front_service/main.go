package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	frontsvc "github.com/mlb/mlb-ballpark-segregation-service/front_service"
	design "github.com/mlb/mlb-ballpark-segregation-service/front_service/design"
	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", "8080", "HTTP port (overrides host HTTP port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	// Setup logging
	var (
		logger *log.Logger
	)
	{
		logger = log.New(os.Stderr, fmt.Sprintf("[%s] ", design.ServiceName), log.Ltime)
	}

	// Initialize the services.
	var (
		schedulerSvc scheduler.Service
	)
	{
		schedulerSvc = frontsvc.NewScheduler(logger)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		schedulerEndpoints *scheduler.Endpoints
	)
	{
		schedulerEndpoints = scheduler.NewEndpoints(schedulerSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":
		{
			addr := "http://localhost:80"
			u, err := url.Parse(addr)
			if err != nil {
				logger.Fatalf("invalid URL %#v: %s\n", addr, err)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h, _, err := net.SplitHostPort(u.Host)
				if err != nil {
					logger.Fatalf("invalid URL %#v: %s\n", u.Host, err)
				}
				u.Host = net.JoinHostPort(h, *httpPortF)
			} else if u.Port() == "" {
				u.Host = net.JoinHostPort(u.Host, "80")
			}
			handleHTTPServer(ctx, u, schedulerEndpoints, &wg, errc, logger, *dbgF)
		}

	default:
		logger.Fatalf("invalid host argument: %q (valid hosts: localhost)\n", *hostF)
	}

	// Wait for signal.
	logger.Printf("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Println("exited")
}

func loadFlagsFromEnv() {
	envToFlags := map[string]string{
		// HTTP Server Arguments
		"HTTPPORT": "http-port",
		"HOST":     "host",
		"DOMAIN":   "domain",
		"DEBUG":    "debug",
	}

	for k, v := range envToFlags {
		val := os.Getenv(k)
		if val == "" {
			continue
		}
		// If supplied, the environment variables will trump the command line arguments.
		os.Args = append(os.Args, fmt.Sprintf("--%s=%s", v, val))
	}
}
