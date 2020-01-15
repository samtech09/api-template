package web

import (
	"log"
	c "github.com/samtech09/api-template/config"
	g "github.com/samtech09/api-template/global"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type key int

const (
	//KeyAppContext is key to get AppContext in middleware
	KeyAppContext key = iota
	//KeyAPIVersion is key to get set API version in context
	KeyAPIVersion key = iota
)

//Server is web server
type Server struct {
	*http.Server
	Router *chi.Mux
}

//NewServer Create new Server object
func NewServer() *Server {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	if c.AppConfig.DisableGzip == false {
		r.Use(middleware.DefaultCompress)
	}

	// Log all requests to this file
	f, _ := os.Create("logs/app.log")
	//DefaultLogger = RequestLogger(&DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: false})
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(f, "", log.LstdFlags), NoColor: true})

	// create http server with custom config
	httpSrv := http.Server{
		Addr:         ":" + strconv.Itoa(c.AppConfig.ListenPort),
		Handler:      r, // < here Chi is attached to the HTTP server
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  1 * time.Minute,
		//  MaxHeaderBytes: 1 << 20,
	}

	//disable keepAlive
	httpSrv.SetKeepAlivesEnabled(false)
	s := Server{&httpSrv, r}

	return &s
}

//Start starts listening the server of set ports
func (s *Server) Start() {
	//ip := GetOutboundIP(c.AppConfig.PingIP)
	listnAddr := "[::1],127.0.0.1," + g.MyIP.String()
	log.Print("Listening on IPs: ", listnAddr)

	addr := strings.Split(listnAddr, ",")
	for _, a := range addr {
		go s.listenAndServeEx(a)
	}
	select {}
}

func (s *Server) listenAndServeEx(addr string) {
	network := "tcp4"
	addr = addr + ":" + strconv.Itoa(c.AppConfig.ListenPort)

	// if addres is enclosed in square brackets then it is IPv6 address
	if strings.HasPrefix(addr, "[") {
		network = "tcp6"
	}
	l, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("ListenEx: ", err)
	}

	if c.AppConfig.DisableSSL {
		// start serving
		err := s.Serve(l)
		if err != nil {
			log.Fatal("ServeEx: ", err)
		}
	} else {
		// start serving with TLS
		err := s.ServeTLS(l, c.AppConfig.SSLCertFile, c.AppConfig.SSLKeyFile)
		if err != nil {
			log.Fatal("ServeExTLS: ", err)
		}
	}
	log.Printf("Listening on %s %s\n", network, addr)
}
