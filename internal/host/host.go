package host

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/internal/protocol"
	"gopkg.in/yaml.v3"
)

//go:embed web/**
var web embed.FS

type hostConfig struct {
	ServicePort int `yaml:"service_port"`
	WebPort     int `yaml:"web_port"`
	ApiPort     int `yaml:"api_port"`
}

func getHost() (*host, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	hostConfigFile, err := os.ReadFile(filepath.Join(cwd, "host-config.yaml"))
	if err != nil {
		return nil, err
	}

	host := host{}
	err = yaml.Unmarshal(hostConfigFile, &host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func Run(ctx context.Context) error {
	host, err := getHost()
	if err != nil {
		return err
	}

	ctx = logging.ContextWithLogger(ctx)
	if err := host.start(ctx); err != nil {
		return err
	}

	// Block forever
	select {}
}

type host struct {
	Config hostConfig `yaml:"host"`
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

	logger.Printf("Starting service listener on port %v\n", h.Config.ServicePort)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", h.Config.ServicePort))
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
		logger.Printf("Starting web server on port %v\n", h.Config.WebPort)

		err = http.ListenAndServe(fmt.Sprintf(":%v", h.Config.WebPort), webMux)
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
		logger.Printf("Starting api server on port %v\n", h.Config.ApiPort)

		err := http.ListenAndServe(fmt.Sprintf(":%v", h.Config.ApiPort), apiMux)
		if err != nil {
			logger.Print(errors.Join(fmt.Errorf("failed to start api server"), err))
		}
	}()

	return nil
}
