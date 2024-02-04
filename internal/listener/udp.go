package listener

import (
	"net"
	"strconv"
	"time"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/metrics"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

const (
	bufSize = 1024 // Size of the buffer to read UDP packets
)

// CreateUDPListener starts a UDP DNS listener on the specified port.
func CreateUDPListener(config *config.Config, resolver *resolver.Resolver) {

	addr := net.JoinHostPort("", strconv.Itoa(config.Protocols.UDP.Port))
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating udp listener")
	}
	defer conn.Close()

	log.Debug().Str("protocol", "udp").Str("address", addr).Msg("listening")
	handleUDPConnections(conn, resolver)
}

// handleUDPConnections listens for incoming UDP packets and processes them.
func handleUDPConnections(conn net.PacketConn, res *resolver.Resolver) {
	buf := make([]byte, bufSize)

	for {
		n, clientAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Error().Err(err).Str("protocol", "udp").Msg("error reading udp packet")
			continue
		}

		go processUDPQuery(buf[:n], conn, clientAddr, res)
	}
}

// processUDPQuery processes a single UDP query and sends back a response.
func processUDPQuery(query []byte, conn net.PacketConn, addr net.Addr, res *resolver.Resolver) {
	startTime := time.Now()
	var req dns.Msg
	if err := req.Unpack(query); err != nil {
		log.Error().Err(err).Str("protocol", "udp").Msg("error unpacking dns query")
		return
	}

	resp, err := res.Resolve(&req)
	if err != nil {
		log.Error().Err(err).Str("protocol", "udp").Msg("error resolving dns query")
		return
	}

	respBytes, err := resp.Pack()
	if err != nil {
		log.Error().Err(err).Str("protocol", "udp").Msg("error packing dns response")
		return
	}

	if _, err := conn.WriteTo(respBytes, addr); err != nil {
		log.Error().Err(err).Str("protocol", "udp").Msg("error sending dns response")
		return
	}

	metrics.RequestDuration.WithLabelValues("udp").Observe(time.Since(startTime).Seconds())
}
