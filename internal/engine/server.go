package engine

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/miekg/dns"

	"github.com/dingdayu/dnsx/pkg/log"
)

type Server struct {
	udp53 *dns.Server
	tcp53 *dns.Server

	httpAPI *http.Server
	httpDNS *http.Server

	logger *log.Logger
}

// NewServer new listen server
func NewServer(opts ...OptionsFunc) *Server {
	// load option
	option := &Option{}
	for _, opt := range opts {
		opt(option) //opt是个方法，入参是*Client，内部会修改client的值
	}

	// new Server struct
	srv := &Server{}

	// new udp server
	if len(option.udpAddr) > 0 && option.dnsHandle != nil {
		srv.udp53 = &dns.Server{Addr: option.udpAddr, Net: "udp", Handler: option.dnsHandle}
	}

	// new tcp server
	if len(option.tcpAddr) > 0 && option.dnsHandle != nil {
		srv.tcp53 = &dns.Server{Addr: option.tcpAddr, Net: "tcp", Handler: option.dnsHandle}
	}

	// new http api server
	if len(option.httpAPIAddr) > 0 && option.httpAPIHandle != nil {
		srv.httpAPI = &http.Server{Addr: option.httpAPIAddr, Handler: option.httpAPIHandle, MaxHeaderBytes: 1 << 20}
	}

	// new http dns server
	if len(option.httpDNSAddr) > 0 && option.httpDNSHandle != nil {
		srv.httpDNS = &http.Server{Addr: option.httpDNSAddr, Handler: option.httpAPIHandle, MaxHeaderBytes: 1 << 20}
	}
	return srv
}

// ListenAndServe listen all serve
func (s Server) ListenAndServe() (err error) {
	err = s.listenTCP53()
	err = s.listenUDP53()
	err = s.listenHttpAPI()
	err = s.listenHttpDNS()
	return err
}

// listenTCP53 listen tcp dns
func (s Server) listenTCP53() (err error) {
	if s.tcp53 != nil {
		go func() {
			err = s.tcp53.ListenAndServe()
			if err != nil {
				fmt.Printf("\u001B[1;30;42m[info]\u001B[0m Start UDP listener on %s failed:%s\n", s.tcp53.Addr, err.Error())
				os.Exit(1)
			}
		}()
	}
	return err
}

// listenUDP53 listen udp dns
func (s Server) listenUDP53() (err error) {
	if s.udp53 != nil {
		go func() {
			err = s.udp53.ListenAndServe()
			if err != nil {
				fmt.Printf("\u001B[1;30;42m[info]\u001B[0m Start TCP listener on %s failed:%s\n", s.udp53.Addr, err.Error())
				os.Exit(1)
			}
		}()
	}
	return err
}

// listenHTTPDNS listen http dns
func (s Server) listenHTTPDNS() (err error) {
	if s.httpDNS != nil {
		go func() {
			if err = s.httpDNS.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Println("\033[1;30;41m[error]\033[0m Start http server error: ", err.Error())
				os.Exit(1)
			}
		}()
	}
	return err
}

// listenHTTPAPI listen http api
func (s Server) listenHTTPAPI() (err error) {
	if s.httpAPI != nil {
		go func() {
			if err = s.httpAPI.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Println("\033[1;30;41m[error]\033[0m Start http server error: ", err.Error())
				os.Exit(1)
			}
		}()
	}
	return err
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s Server) Shutdown(ctx context.Context) (err error) {
	// HTTP Shutdown
	if s.httpAPI != nil {
		if err = s.httpAPI.Shutdown(ctx); err != nil {
			fmt.Printf("\033[1;30;43m[warn]\033[0m http api Shutdown: %s\n", err)
		}
	}

	if s.httpDNS != nil {
		if err = s.httpDNS.Shutdown(ctx); err != nil {
			fmt.Printf("\033[1;30;43m[warn]\033[0m http dns Shutdown: %s\n", err)
		}
	}

	// DNS Shutdown
	if s.udp53 != nil {
		if err = s.udp53.ShutdownContext(ctx); err != nil {
			fmt.Printf("\033[1;30;43m[warn]\033[0m UDP Shutdown: %s\n", err)
		}

	}
	if s.tcp53 != nil {
		if err = s.tcp53.ShutdownContext(ctx); err != nil {
			fmt.Printf("\033[1;30;43m[warn]\033[0m TCP Shutdown: %s\n", err)
		}
	}
	return err
}
