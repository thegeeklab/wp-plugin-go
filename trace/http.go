// Copyright (c) 2019, Drone Plugins project authors
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package trace

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http/httptrace"
	"net/textproto"

	"github.com/rs/zerolog/log"
)

// HTTP uses httptrace to log all network activity for HTTP requests.
func HTTP(ctx context.Context) context.Context {
	return httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.Trace().Str("host-port", hostPort).Msg("ClientTrace.GetConn")
		},

		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Trace().
				Str("local-address", connInfo.Conn.LocalAddr().String()).
				Str("remote-address", connInfo.Conn.RemoteAddr().String()).
				Bool("reused", connInfo.Reused).
				Bool("was-idle", connInfo.WasIdle).
				Dur("idle-time", connInfo.IdleTime).
				Msg("ClientTrace.GoConn")
		},

		PutIdleConn: func(err error) {
			log.Trace().Err(err).Msg("ClientTrace.GoConn")
		},

		GotFirstResponseByte: func() {
			log.Trace().Msg("ClientTrace.GotFirstResponseByte")
		},

		Got100Continue: func() {
			log.Trace().Msg("ClientTrace.Got100Continue")
		},

		Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
			log.Trace().
				Int("code", code).
				Str("header", fmt.Sprint(header)).
				Msg("ClientTrace.Got1xxxResponse")

			return nil
		},

		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			log.Trace().Str("host", dnsInfo.Host).Msg("ClientTrace.DNSStart")
		},

		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			log.Trace().
				Str("addresses", fmt.Sprint(dnsInfo.Addrs)).
				Err(dnsInfo.Err).
				Bool("coalesced", dnsInfo.Coalesced).
				Msg("ClientTrace.DNSDone")
		},

		ConnectStart: func(network, addr string) {
			log.Trace().
				Str("network", network).
				Str("address", addr).
				Msg("ClientTrace.ConnectStart")
		},

		ConnectDone: func(network, addr string, err error) {
			log.Trace().
				Str("network", network).
				Str("address", addr).
				Err(err).
				Msg("ClientTrace.ConnectDone")
		},

		TLSHandshakeStart: func() {
			log.Trace().Msg("ClientTrace.TLSHandshakeStart")
		},

		TLSHandshakeDone: func(connState tls.ConnectionState, err error) {
			log.Trace().
				Uint16("version", connState.Version).
				Bool("handshake-complete", connState.HandshakeComplete).
				Bool("did-resume", connState.DidResume).
				Uint16("cipher-suite", connState.CipherSuite).
				Str("negotiated-protocol", connState.NegotiatedProtocol).
				Str("server-name", connState.ServerName).
				Err(err).
				Msg("ClientTrace.TLSHandshakeDone")
		},

		WroteHeaderField: func(key string, value []string) {
			log.Trace().
				Str("key", key).
				Strs("values", value).
				Msg("ClientTrace.WroteHeaderField")
		},

		WroteHeaders: func() {
			log.Trace().Msg("ClientTrace.WroteHeaders")
		},

		Wait100Continue: func() {
			log.Trace().Msg("ClientTrace.Wait100Continue")
		},

		WroteRequest: func(reqInfo httptrace.WroteRequestInfo) {
			log.Trace().Err(reqInfo.Err).Msg("ClientTrace.WroteRequest")
		},
	})
}
