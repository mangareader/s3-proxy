package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"emperror.dev/errors"
	"github.com/dimiro1/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/config"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/log"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/metrics"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/server/middlewares"
	"github.com/oxyno-zeta/s3-proxy/pkg/s3-proxy/tracing"
)

type InternalServer struct {
	logger     log.Logger
	cfgManager config.Manager
	metricsCl  metrics.Client
	server     *http.Server
}

func NewInternalServer(logger log.Logger, cfgManager config.Manager, metricsCl metrics.Client) *InternalServer {
	return &InternalServer{
		logger:     logger,
		cfgManager: cfgManager,
		metricsCl:  metricsCl,
	}
}

func (svr *InternalServer) Listen() error {
	svr.logger.Infof("Internal server listening on %s", svr.server.Addr)

	var err error

	// Listen (either HTTPS or HTTP)
	if svr.server.TLSConfig != nil {
		err = svr.server.ListenAndServeTLS("", "")
	} else {
		err = svr.server.ListenAndServe()
	}

	// Check error
	if err != nil {
		return errors.WithStack(err)
	}

	// Default
	return nil
}

func (svr *InternalServer) GenerateServer() error {
	// Get configuration
	cfg := svr.cfgManager.GetConfig()
	// Generate internal router
	r := svr.generateInternalRouter()
	// Create server
	addr := cfg.InternalServer.ListenAddr + ":" + strconv.Itoa(cfg.InternalServer.Port)
	server := &http.Server{ //nolint: gosec // Set after
		Addr:    addr,
		Handler: r,
	}

	// Inject timeouts
	err := injectServerTimeout(server, cfg.InternalServer.Timeouts)
	// Check error
	if err != nil {
		return err
	}

	// Get the TLS configuration (if necessary)
	tlsConfig, err := generateTLSConfig(cfg.InternalServer.SSL, svr.logger)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create TLS configuration for internal server: %v", err))
	}

	server.TLSConfig = tlsConfig

	// Store server
	svr.server = server

	return nil
}

func (svr *InternalServer) generateInternalRouter() http.Handler {
	r := chi.NewRouter()

	// Get configuration
	cfg := svr.cfgManager.GetConfig()

	// Check if we need to enabled the compress middleware
	if *cfg.InternalServer.Compress.Enabled {
		r.Use(middleware.Compress(
			cfg.InternalServer.Compress.Level,
			cfg.InternalServer.Compress.Types...,
		))
	}

	// Check if no cache is disabled or not
	if cfg.InternalServer.Cache == nil || cfg.InternalServer.Cache.NoCacheEnabled {
		// Apply no cache
		r.Use(middleware.NoCache)
	} else {
		// Apply S3 proxy cache management middleware
		r.Use(middlewares.CacheManagement(cfg.InternalServer.Cache))
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(log.NewStructuredLogger(
		svr.logger,
		tracing.GetTraceIDFromRequest,
	))
	r.Use(log.HTTPAddLoggerToContextMiddleware())
	r.Use(svr.metricsCl.Instrument("internal", cfg.Metrics))
	r.Use(middleware.Recoverer)

	healthHandler := health.NewHandler()
	// Listen path
	r.Handle("/metrics", svr.metricsCl.GetExposeHandler())
	r.Handle("/health", healthHandler)
	r.Handle("/config", configHandler(svr.cfgManager))

	return r
}

func configHandler(cfgManager config.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Get configuration
		cfg := cfgManager.GetConfig()

		// Create output answer
		type resp struct {
			Config *config.Config `json:"config"`
		}
		// json marshal
		bb, err := json.Marshal(&resp{Config: cfg})
		// Check error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))

			// Stop
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		_, _ = w.Write(bb)
	})
}
