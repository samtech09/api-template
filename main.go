package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/samtech09/api-template/controllers"
	g "github.com/samtech09/api-template/global"
	"github.com/samtech09/api-template/internal/initializer"
	"github.com/samtech09/api-template/web"

	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	appVer   string
	buildVer string
)

func main() {
	//parse commandline flags
	prod := flag.Bool("prod", false, "-prod to run in production mode")
	versionFlag := flag.Bool("v", false, "Print the current version and exit")
	flag.Parse()

	switch {
	case *versionFlag:
		printVersion()
		return
	}

	//
	//initialize config
	//
	initializer.Initconfig(prod, "")

	//
	//get path to current folder
	//
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	//
	//initialize logs
	//
	initializer.InitLogger(dir)

	//
	//initialize services and connections
	//
	initializer.InitServices(dir)

	// setup signal catching
	sigs := make(chan os.Signal, 1)

	//
	// catch all signals since not explicitly listing
	//
	signal.Notify(sigs)
	// method invoked upon seeing signal
	go func() {
		s := <-sigs
		if s == syscall.SIGPIPE {
			g.Logger.Info().Str("SIGNAL", s.String()).Msg("RECEIVED [IGNORED] SIGNAL")
		} else {
			g.Logger.Info().Str("SIGNAL", s.String()).Msg("RECEIVED TERMINATE SIGNAL")
			initializer.AppCleanup()
			os.Exit(1)
		}
	}()

	g.Logger.Info().Str("AppVersion", appVer).Str("Build", buildVer).Msg("Current Version and build")
	if g.IsProduction {
		g.Logger.Info().Msg("Server started in **production** mode")
	} else {
		g.Logger.Info().Msg("Server started in development mode")
	}
	g.Logger.Info().Msgf("Running for: %s", g.Config.MyURL)
	g.Logger.Info().Msgf("LogLevel: %s", g.Config.LogLevel)

	s := web.NewServer()
	v1 := web.NewAPIVersion("v1", s.Router)

	s.Router.Get("/", index)

	// Register controllers
	controllers.Register(&controllers.User{}, v1)

	defer initializer.AppCleanup()
	s.Start()
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to API server from API-Template")
}

func printVersion() {
	log.Printf("Current version and build: %s %s", appVer, buildVer)
}
