package main

import (
	"log"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/listener"
	"github.com/bwoff11/go-resolve/internal/models"
	"github.com/bwoff11/go-resolve/internal/resolver"
)

func main() {
	config.Load()

	res := resolver.New()

	tcpListener := listener.New(res, models.TCP, 1053)
	udpListener := listener.New(res, models.UDP, 1053)

	go func() {
		if err := tcpListener.Listen(); err != nil {
			log.Fatalf("TCP Listener failed: %v", err)
		}
	}()

	go func() {
		if err := udpListener.Listen(); err != nil {
			log.Fatalf("UDP Listener failed: %v", err)
		}
	}()

	select {}
}
