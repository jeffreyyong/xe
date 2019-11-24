package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type XEService struct {
	*http.Server
}

func NewXEService(h http.Handler, addr string) *XEService {
	server := &http.Server{
		Handler: h,
		Addr:    addr,
	}

	return &XEService{
		Server: server,
	}
}

func (x *XEService) Run() {
	if err := x.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("XE service failed to start up: ", err)
	}
}

func (x *XEService) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := x.Shutdown(ctx); err != nil {
		log.Fatal("XE service failed to shut down: ", err)
	}
}
