package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
)

/**
 * @Author: v_rainlliu
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2020/12/9 15:47
 */


func NewServer(ctx context.Context,addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})
	server := &http.Server{
		Addr: addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Println("server stop")
		server.Shutdown(ctx)
	}()
	return server.ListenAndServe()
}

func main()  {
	g,ctx := errgroup.WithContext(context.Background())

	//http服务1
	g.Go(func() error {
		if err := NewServer(ctx,"8080"); err != nil {
			return err
		}
		return nil
	})
	//http服务2
	g.Go(func() error {
		if err := NewServer(ctx,"8081");err != nil {
			return err
		}
		return nil
	})
	//信号监听
	g.Go(func() error {
		s := make(chan os.Signal)
		signal.Notify(s)
		select {
		case <-s:
			return errors.New("receive exit signal")
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != nil {
		log.Println("服务已退出")
	}
}
