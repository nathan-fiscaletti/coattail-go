package host

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/nathan-fiscaletti/coattail-go/internal/host/api"
	"github.com/nathan-fiscaletti/coattail-go/internal/host/config"
	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/internal/logging"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailtypes"
)

var (
	ErrHostNotFound = errors.New("host not found in context")
	ErrInvalidHost  = errors.New("invalid host found in context")
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
		return nil, ErrHostNotFound
	}

	if h, ok := host.(*Host); ok {
		return h, nil
	}

	return nil, ErrInvalidHost
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

type ConnectionHandler func(context.Context, net.Conn, bool)

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

func (h *Host) startListener(ctx context.Context, handleConnection ConnectionHandler) error {
	certFile := "server.crt"
	keyFile := "server.key"

	_, certFileErr := os.Stat(certFile)
	_, keyFileErr := os.Stat(keyFile)

	certFileExists := certFileErr == nil || !os.IsNotExist(certFileErr)
	keyFileExists := keyFileErr == nil || !os.IsNotExist(keyFileErr)

	if !certFileExists || !keyFileExists {
		if logger, err := logging.GetLogger(ctx); err == nil {
			logger.Printf("certificate or key missing, generating self-signed certificate\n")
		}

		// delete cert and key file if they exist

		if certFileErr == nil {
			err := os.Remove(certFile)
			if err != nil {
				return fmt.Errorf("failed to delete existing certificate file: %w", err)
			}
		}

		if keyFileErr == nil {
			err := os.Remove(keyFile)
			if err != nil {
				return fmt.Errorf("failed to delete existing key file: %w", err)
			}
		}

		// generate new cert and key

		err := h.createSelfSignedCertificate(ctx, h.Config.ServiceConfig.Address.Host)
		if err != nil {
			return fmt.Errorf("failed to generate self-signed certificate: %w", err)
		}
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load certificate and key: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := net.Listen("tcp", h.Config.ServiceConfig.Address.String())
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	tlsListener := tls.NewListener(listener, tlsConfig)
	go func() {
		defer listener.Close()

		for {
			conn, err := tlsListener.Accept()
			if err != nil {
				if logger, _ := logging.GetLogger(ctx); logger != nil {
					logger.Println(fmt.Errorf("failed to accept connection: %v", err))
				}
				break
				// continue
			}

			go handleConnection(ctx, conn, h.Config.ServiceConfig.LogPackets)
		}
	}()

	if logger, _ := logging.GetLogger(ctx); logger != nil {
		logger.Printf("running service at %s\n", h.Config.ServiceConfig.Address.String())
	}

	return nil
}

func (h *Host) startWebServer(ctx context.Context) error {
	if !h.Config.WebConfig.Enabled {
		return nil
	}

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

		if logger, err := logging.GetLogger(ctx); err == nil {
			logger.Printf("running web server at %v\n", h.Config.WebConfig.Address.String())
		}

		err = http.ListenAndServe(h.Config.WebConfig.Address.String(), webMux)
		if err != nil {
			if logger, err := logging.GetLogger(ctx); err == nil {
				logger.Print(fmt.Errorf("listen and serve error: %w", err))
			}
		}
	}()

	return nil
}

func loggingMiddleware(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger, err := logging.GetLogger(ctx); err == nil {
			logger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Host) startApiServer(ctx context.Context) error {
	if !h.Config.ApiConfig.Enabled {
		return nil
	}

	go func() {
		if logger, err := logging.GetLogger(ctx); err == nil {
			apiLogger := log.New(os.Stdout, logger.Prefix()+"[API] ", log.LstdFlags)
			ctx = context.WithValue(ctx, keys.LoggerKey, apiLogger)
		}

		apiMux := http.NewServeMux()

		apiMux.Handle("/healthcheck", loggingMiddleware(ctx, api.NewHealthCheckHandler(ctx, h.LocalPeer)))
		apiMux.Handle("/peers", loggingMiddleware(ctx, api.NewPeersHandler(ctx, h.LocalPeer)))
		apiMux.Handle("/actions", loggingMiddleware(ctx, api.NewActionsHandler(ctx, h.LocalPeer)))

		if logger, err := logging.GetLogger(ctx); err == nil {
			logger.Printf("running api server at %v\n", h.Config.ApiConfig.Address.String())
		}

		err := http.ListenAndServe(h.Config.ApiConfig.Address.String(), apiMux)
		if err != nil {
			if logger, err := logging.GetLogger(ctx); err == nil {
				logger.Print(fmt.Errorf("failed to start api server: %w", err))
			}
		}
	}()

	return nil
}
