package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dingdayu/dnsx/api/router"
	"github.com/dingdayu/dnsx/internal/engine"
	"github.com/dingdayu/dnsx/pkg/config"
)

// Run server command
func Run() {
	// dns addr
	dnsAddr := fmt.Sprintf(":%d", config.GetInt("app.port"))
	apiAddr := config.GetString("api.addr")

	// init router
	router.Init()

	// get handler
	dh := router.DNSHandler()
	hh := router.HTTPHandler()

	srv := engine.NewServer(
		engine.WithDNSHandle(dh),
		engine.WithHttpAPIHandle(hh),

		engine.WithDNSAddr(dnsAddr),
		engine.WithHttpAPIAddr(apiAddr),
	)

	// listen serve
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("\033[1;30;41m[error]\033[0m Start server error: ", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("\033[1;30;42m[info]\033[0m Start server listening")
	}

	// Safe exit via signal
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("\n\033[1;30;42m[info]\033[0m Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// server shutdown
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("\033[1;30;43m[warn]\033[0m Server Shutdown: %s\n", err)
	}
	fmt.Println("Server exited")
}
