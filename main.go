package main

import (
	"context"
	"errors"
	"flag"
	"fs_monitor/conf"
	"fs_monitor/goesl"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var rc = flag.String("conf", "./fmrc", "Specify the location of fs_monitor")

type running struct {
	ctx    context.Context
	cancel func()
	conf   *conf.Conf
	wg     *sync.WaitGroup
	conns  []*goesl.Connection
}

var G running

func main() {
	confFile, err := os.OpenFile(*rc, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("open configuration file: %v", err)
	}
	G.conf, err = conf.ConfParse(confFile)
	if err != nil {
		log.Fatalf("configuration parse: %v", err)
	}
	G.wg = &sync.WaitGroup{}
	handleInstances(G)
	G.ctx, G.cancel = context.WithCancel(context.Background())
	go sigCatcher(G)
	G.wg.Wait() // wait all go routine finished
}

func handleInstances(r running) {
	for _, conn := range r.conns {
		r.wg.Add(1)
		go func(ctx context.Context, c *goesl.Connection, wg *sync.WaitGroup) {
			err := c.HandleEvents(ctx)
			if errors.Is(err, net.ErrClosed) || errors.Is(err, context.Canceled) {
				goesl.Noticef("process exiting...")
			} else {
				goesl.Errorf("exiting with error: %v", err)
			}
			wg.Done()
		}(r.ctx, conn, r.wg)
	}
}

func sigCatcher(r running) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigs)
	sig := <-sigs
	goesl.Noticef("got signal %v, quiting ...", sig)
	r.cancel()
	for _, conn := range r.conns {
		conn.Close()
	}
}
