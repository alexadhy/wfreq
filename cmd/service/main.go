package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/alexadhy/wfreq/internal/api"
	"github.com/alexadhy/wfreq/internal/logging"
)

func signalContext(ctx context.Context, log logrus.FieldLogger) (context.Context, context.CancelFunc) {
	if log == nil {
		log = logrus.New()
	}

	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	listenSigs := listenSignals()
	signal.Notify(ch, listenSigs...)

	go func() {
		select {
		case sig := <-ch:
			log.WithField("signal", sig).
				Info("Closing with received signal.")
		case <-ctx.Done():
		}
		cancel()
	}()

	return ctx, cancel
}

func listenSignals() []os.Signal {
	return []os.Signal{syscall.SIGTERM}
}

func main() {
	l := logging.New()
	var addr string
	var port int
	flag.StringVar(&addr, "a", "127.0.0.1", "address to listen to (default to localhost)")
	flag.IntVar(&port, "p", 3334, "port to use (default 3334)")
	flag.Parse()

	wfreqSvc := api.New(*l)
	addrPort := fmt.Sprintf("%s:%d", addr, port)
	l.Infof("Listening on %s", addrPort)
	ctx, cancel := signalContext(context.Background(), l)
	if err := http.ListenAndServe(":3334", wfreqSvc); err != nil {
		l.Fatalf("Error: %v", err)
		cancel()
	}
	<-ctx.Done()
}
