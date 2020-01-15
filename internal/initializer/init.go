package initializer

import (
	"log"
	"net"
	"os"
	"path/filepath"

	c "github.com/samtech09/api-template/config"
	g "github.com/samtech09/api-template/global"
	"github.com/samtech09/api-template/internal/logger"
	"github.com/samtech09/api-template/psql"
	"github.com/samtech09/api-template/sqls"
	"github.com/samtech09/api-template/viewmodels"
	"github.com/samtech09/apiclient"
	"github.com/samtech09/apiroutecache"
	"github.com/samtech09/jwtauth"
	"github.com/samtech09/redicache"
)

var logFile *os.File

// InitServices initializes connections / sessions with different services
func InitServices(appDir string) {
	// get IP address of outbound interface
	g.MyIP = getOutboundIP(c.AppConfig.PingIP)

	//
	// inititalize route cache
	//
	g.Mgosesion = apiroutecache.InitSession(c.AppConfig.Mongo)

	//
	// inititalize data cache
	//
	g.Cache = redicache.InitSession(c.AppConfig.Redis)

	//
	// initialize JWT Validator
	//

	g.JWTval = jwtauth.InitValidator(c.AppConfig.JwtAuthConfig, appDir+"/public.pem")

	//
	//initialize Api
	//
	//g.Api = apiclient.API{AllowInsecureSSL: !c.IsProduction, Timeout: 10}
	g.Api = apiclient.API{AllowInsecureSSL: false, Timeout: 10}

	// ----------------------
	// Initialize db connections
	//
	g.Db = psql.InitDbPool(c.AppConfig.DBReader, c.AppConfig.DBWriter, *g.Logger)

	//
	// setup models for caching
	//
	viewmodels.RegisterModelsForCaching()

	//
	// Load SQLs from JSON
	//
	sqls.LoadSQLs(appDir + "/sqls/sqlbuilder.json")
}

//InitLogger initialize logging to file
func InitLogger(appDir string) {
	loginfo := true
	logerror := true
	logdebug := false

	switch c.AppConfig.LogLevel {
	case "fatal":
		logerror = false
		loginfo = false
		logdebug = false
	case "errors":
		logerror = true
		loginfo = false
		logdebug = false
	case "info":
		logerror = true
		loginfo = true
		logdebug = false
	case "debug":
		logerror = true
		loginfo = true
		logdebug = true
	}

	g.Logger = logger.New(logerror, loginfo, logdebug)

	// set output to file
	//create logs directory if not exist
	newpath := filepath.Join(appDir, "logs")
	os.Mkdir(newpath, os.ModePerm)
	var err error
	logFile, err = os.Open(filepath.Join(newpath, "app.log"))
	if err != nil {
		g.Logger.Fatal().Err(err).Msg("Log file open error")
	}

	g.Logger.Output(logFile)
}

//AppCleanup cleanup existing connections and sessions
func AppCleanup() {
	g.Mgosesion.Cleanup()
	g.Logger.Info().Msg("Shutting down...")
	logFile.Sync()
	logFile.Close()
}

// Get preferred outbound ip of this machine
func getOutboundIP(pingIP string) net.IP {
	//conn, err := net.Dial("udp", "8.8.8.8:80")
	conn, err := net.Dial("udp", pingIP)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
