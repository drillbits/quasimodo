package main

import (
	"flag"
	"fmt"
	qsmd "github.com/drillbits/quasimodo/quasimodo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configFile := flag.String("conf", qsmd.DefaultConfig.ConfigFile, "Config file path")
	flag.Parse()

	conf, err := qsmd.LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	qsmd.NewStore(conf)

	h := qsmd.NewHTTPHandler(conf)
	go http.ListenAndServe(fmt.Sprintf("%s:%s", conf.Host, conf.Port), h)

	w := qsmd.NewWatcher(conf)
	go w.Watch()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	exit := make(chan int)
	go func() {
		for {
			s := <-sig
			switch s {
			case syscall.SIGHUP:
				// kill -SIGHUP pid
				log.Println("SIGHUP")
			exit <- 0
			case syscall.SIGINT:
				// kill -SIGINT pid or Ctrl+c
				log.Println("SIGINT")
			exit <- 0
			case syscall.SIGTERM:
				// kill -SIGTERM pid
				log.Println("SIGTERM")
			exit <- 0
			case syscall.SIGQUIT:
				// kill -SIGQUIT pid
				log.Println("SIGQUIT")
			exit <- 0
			default:
				log.Println("unknown signal")
			exit <- 1
			}
		}
	}()

	code := <-exit
	os.Exit(code)
}
