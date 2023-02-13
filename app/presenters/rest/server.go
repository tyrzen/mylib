package rest

import (
	"crypto/tls"
	"fmt"

	"golang.org/x/crypto/acme/autocert"

	"net/http"
	"os"
	"time"
)

type Server struct{ *http.Server }

func NewServer(hdl http.Handler) (*Server, error) {
	addr := os.Getenv("SRV_HOST") + ":" + os.Getenv("SRV_PORT")

	readTimeout, err := time.ParseDuration(os.Getenv("SRV_READ_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("failed parse read timeout: %w", err)
	}

	writeTimeout, err := time.ParseDuration(os.Getenv("SRV_WRITE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("failed parse write timeout: %w", err)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("SRV_IDLE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("failed parse idle timeout: %w", err)
	}

	dirCache := os.Getenv("SRV_CERT_DIR")
	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(dirCache),
	}
	tlsConfig := &tls.Config{GetCertificate: certManager.GetCertificate}

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      hdl,
		TLSConfig:    tlsConfig,
	}

	return &Server{srv}, nil
}

func (srv *Server) Run() error {
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("error running the server: %w", err)
	}

	if err := srv.ListenAndServeTLS("", ""); err != nil {
		return fmt.Errorf("error running the tls server: %w", err)
	}

	return nil
}
