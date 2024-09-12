package host

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"

	"github.com/nathan-fiscaletti/coattail-go/internal/host/config"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/internal/protocol"
)

//go:embed web/**
var web embed.FS

var hostValue *host

func GetHost() (*host, error) {
	var err error
	if hostValue == nil {
		hostValue, err = newHost()
	}

	return hostValue, err
}

func newHost() (*host, error) {
	hostConfig, err := config.GetHostConfig()
	if err != nil {
		return nil, err
	}

	return &host{
		Config: hostConfig,
	}, nil
}

func Run(ctx context.Context, after func()) error {
	host, err := GetHost()
	if err != nil {
		return err
	}

	ctx = logging.ContextWithLogger(ctx)
	if err := host.start(ctx); err != nil {
		return err
	}

	if after != nil {
		after()
	}

	// Block forever
	select {}
}

type host struct {
	Config *config.HostConfig `yaml:"host"`
}

func (h *host) start(ctx context.Context) error {
	var err error

	if err = h.startListener(ctx); err == nil {
		if err = h.startApiServer(ctx); err == nil {
			if err = h.startWebServer(ctx); err == nil {
				return nil
			}
		}
	}

	return err
}

func (h *host) startListener(ctx context.Context) error {
	logger := logging.GetLogger(ctx)

	logger.Printf("starting service at %v\n", h.Config.ServiceAddress)
	listener, err := net.Listen("tcp", h.Config.ServiceAddress)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to start listener"), err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Print(errors.Join(fmt.Errorf("failed to accept connection"), err))
				continue
			}

			go func() {
				go protocol.NewPacketHandler(ctx, conn).HandlePackets()
			}()
		}
	}()

	return nil
}

func (h *host) startWebServer(ctx context.Context) error {
	logger := logging.GetLogger(ctx)

	wfbFs, err := fs.Sub(web, "web")
	if err != nil {
		// We should panic here because this means the embedded
		// filesystem is not working correctly.
		panic(err)
	}

	go func() {
		webMux := http.NewServeMux()
		fs := http.FileServer(http.FS(wfbFs))
		webMux.Handle("/", fs)
		logger.Printf("starting web server at %v\n", h.Config.WebAddress)

		err = http.ListenAndServe(h.Config.WebAddress, webMux)
		if err != nil {
			logger.Print(errors.Join(fmt.Errorf("listen and serve error"), err))
		}
	}()

	return nil
}

func (h *host) startApiServer(ctx context.Context) error {
	logger := logging.GetLogger(ctx)

	go func() {
		apiMux := http.NewServeMux()
		apiMux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})
		logger.Printf("starting api server at %v\n", h.Config.ApiAddress)

		err := http.ListenAndServe(h.Config.ApiAddress, apiMux)
		if err != nil {
			logger.Print(errors.Join(fmt.Errorf("failed to start api server"), err))
		}
	}()

	return nil
}
