package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"dnsx/api/router"
	"dnsx/pkg/config"
	"dnsx/pkg/validate"

	"github.com/gin-gonic/gin/binding"
	"github.com/miekg/dns"
)

func Run() {
	// 注册 DNS
	dnsAddr := fmt.Sprintf(":%d", config.GetInt("app.port"))
	udp := &dns.Server{Addr: dnsAddr, Net: "udp", Handler: &DNSHandler{}}
	go func() {
		err := udp.ListenAndServe()
		if err != nil {
			log.Fatalf("\u001B[1;30;42m[info]\u001B[0m Start UDP listener on %s failed:%s\n", dnsAddr, err.Error())
		}
	}()

	tcp := &dns.Server{Addr: dnsAddr, Net: "tcp", Handler: &DNSHandler{}}
	go func() {
		err := tcp.ListenAndServe()
		if err != nil {
			log.Fatalf("\u001B[1;30;42m[info]\u001B[0m Start TCP listener on %s failed:%s\n", dnsAddr, err.Error())
		}
	}()

	// 注册 HTTP
	addr := config.GetString("api.addr")
	binding.Validator = validate.GinValidator()
	srv := &http.Server{
		Addr:           addr,
		Handler:        router.Handler(),
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("\033[1;30;42m[info]\033[0m Start http server listening %s\n", addr)
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("\033[1;30;41m[error]\033[0m Start http server error: ", err.Error())
			os.Exit(1)
		}
	}()

	// Safe exit via signal
	// 5 Second
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("\n\033[1;30;42m[info]\033[0m Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// HTTP Shutdown
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("\033[1;30;42m[info]\033[0m Api Shutdown: %s", err)
	}

	// DNS Shutdown
	if err := udp.ShutdownContext(ctx); err != nil {
		fmt.Printf("\033[1;30;42m[info]\033[0m UDP Shutdown: %s", err)
	}
	if err := tcp.ShutdownContext(ctx); err != nil {
		fmt.Printf("\033[1;30;42m[info]\033[0m TCP Shutdown: %s", err)
	}

	fmt.Println("Server exited")
}
