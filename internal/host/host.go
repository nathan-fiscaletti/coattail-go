package host

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"

	"github.com/nathan-fiscaletti/coattail-go/internal/host/config"
	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

//go:embed web/**
var web embed.FS

type Host struct {
	Config    *config.HostConfig `yaml:"host"`
	LocalPeer *coattailtypes.Peer
}

func ContextWithHost(ctx context.Context) (context.Context, error) {
	host := ctx.Value(keys.HostKey)

	if host == nil {
		host, err := newHost()
		if err != nil {
			return nil, err
		}

		return context.WithValue(ctx, keys.HostKey, host), nil
	}

	return ctx, nil
}

func GetHost(ctx context.Context) (*Host, error) {
	host := ctx.Value(keys.HostKey)
	if host == nil {
		return nil, fmt.Errorf("no host found in context")
	}

	if h, ok := host.(*Host); ok {
		return h, nil
	}

	return nil, fmt.Errorf("invalid host found in context")
}

func newHost() (*Host, error) {
	hostConfig, err := config.GetHostConfig()
	if err != nil {
		return nil, err
	}

	return &Host{
		Config: hostConfig,
	}, nil
}

type ConnectionHandler func(context.Context, net.Conn)

func (h *Host) Start(ctx context.Context, connHandler ConnectionHandler) error {
	var err error

	if err = h.startListener(ctx, connHandler); err == nil {
		if err = h.startApiServer(ctx); err == nil {
			if err = h.startWebServer(ctx); err == nil {
				return nil
			}
		}
	}

	return err
}

func (h *Host) startListener(ctx context.Context, connHandler ConnectionHandler) error {
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
				connHandler(ctx, conn)
			}()
		}
	}()

	return nil
}

func (h *Host) startWebServer(ctx context.Context) error {
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

func (h *Host) startApiServer(ctx context.Context) error {
	logger := logging.GetLogger(ctx)

	go func() {
		apiMux := http.NewServeMux()
		apiMux.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
			peers, err := h.LocalPeer.ListPeers(ctx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			peerData, err := json.Marshal(peers)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			// disable cors
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(peerData)
		})

		logger.Printf("starting api server at %v\n", h.Config.ApiAddress)

		err := http.ListenAndServe(h.Config.ApiAddress, apiMux)
		if err != nil {
			logger.Print(errors.Join(fmt.Errorf("failed to start api server"), err))
		}
	}()

	return nil
}
