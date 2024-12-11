// Copyright 2023 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	plugin_trace "github.com/thegeeklab/wp-plugin-go/v4/trace"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/proxy"
)

const (
	NetDailerTimeout                 = 30 * time.Second
	HTTPTransportIdleTimeout         = 90 * time.Second
	HTTPTransportTLSHandshakeTimeout = 10 * time.Second
	HTTPTransportMaxIdleConns        = 100
)

// Network contains options for connecting to the network.
type Network struct {
	// Context for making network requests.
	//
	// If `trace` logging is requested the context will use `httptrace` to
	// capture all network requests.
	//nolint:containedctx
	Context context.Context

	/// Whether SSL verification is skipped or not.
	InsecureSkipVerify bool

	// Client for making network requests.
	Client *http.Client
}

func networkFlags(category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     "transport.insecure-skip-verify",
			Usage:    "skip SSL verification",
			EnvVars:  []string{"PLUGIN_INSECURE_SKIP_VERIFY"},
			Category: category,
		},
		&cli.StringFlag{
			Name:    "transport.socks-proxy",
			Usage:   "socks proxy address",
			EnvVars: []string{"SOCKS_PROXY"},
			Hidden:  true,
		},
		&cli.BoolFlag{
			Name:    "transport.socks-proxy-off",
			Usage:   "socks proxy ignored",
			EnvVars: []string{"SOCKS_PROXY_OFF"},
			Hidden:  true,
		},
	}
}

func NetworkFromContext(ctx *cli.Context) Network {
	var (
		skipVerify     = ctx.Bool("transport.insecure-skip-verify")
		defaultContext = context.Background()
		socks          = ctx.String("transport.socks-proxy")
		socksoff       = ctx.Bool("transport.socks-proxy-off")
	)

	certs, err := x509.SystemCertPool()
	if err != nil {
		log.Error().Err(err).Msg("failed to find system CA certs")
	}

	tlsConfig := &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: skipVerify, //nolint:gosec
	}

	transport := &http.Transport{
		TLSClientConfig:       tlsConfig,
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          HTTPTransportMaxIdleConns,
		IdleConnTimeout:       HTTPTransportIdleTimeout,
		TLSHandshakeTimeout:   HTTPTransportTLSHandshakeTimeout,
		ExpectContinueTimeout: 1 * time.Second,
	}

	dialer := &net.Dialer{
		Timeout:   NetDailerTimeout,
		KeepAlive: NetDailerTimeout,
		DualStack: true,
	}

	if len(socks) != 0 && !socksoff {
		proxyDialer, err := proxy.SOCKS5("tcp", socks, nil, dialer)
		if err != nil {
			log.Error().Err(err).Msg("failed to create socks proxy")
		}

		if contextDialer, ok := proxyDialer.(proxy.ContextDialer); ok {
			transport.DialContext = contextDialer.DialContext
		} else {
			transport.DialContext = func(_ context.Context, network, addr string) (net.Conn, error) {
				return proxyDialer.Dial(network, addr)
			}
		}
	} else {
		transport.DialContext = dialer.DialContext
	}

	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		defaultContext = plugin_trace.HTTP(defaultContext)
	}

	client := &http.Client{
		Transport: transport,
	}

	return Network{
		Context:            defaultContext,
		InsecureSkipVerify: skipVerify,
		Client:             client,
	}
}
