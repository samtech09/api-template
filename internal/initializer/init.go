package initializer

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"

	g "github.com/samtech09/api-template/global"
	"github.com/samtech09/api-template/internal/logger"
	"github.com/samtech09/api-template/sqls"
	"github.com/samtech09/api-template/viewmodels"
	"github.com/samtech09/apiclient"
	"github.com/samtech09/apiroutecache"
	"github.com/samtech09/dbtools/mango"
	"github.com/samtech09/dbtools/pgsql"
	"github.com/samtech09/jwtauth"
	"github.com/samtech09/redicache"
)

var logFile *os.File

//Initconfig read config file and initialize AppConfig struct
func Initconfig(prod *bool, confFolderPath string) {
	//set if productin mode is true
	g.IsProduction = *prod

	conffile := confFolderPath + "conf.dev.json"
	if g.IsProduction {
		// fmt.Println("****************************************")
		// fmt.Println("*** Loading PRODUCTION configuration ***")
		// fmt.Println("****************************************")
		conffile = confFolderPath + "conf.prod.json"
	}

	file, err := os.Open(conffile)
	if err != nil {
		log.Fatal("Missing config file.\n", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	appConfig := g.NewConfig()
	err = decoder.Decode(&appConfig)
	if err != nil {
		log.Fatal("Cannot parse config.\n", err)
	}

	if g.IsProduction {
		appConfig.MyURL = "https://" + appConfig.MyDomain
	} else {
		if appConfig.DisableSSL {
			appConfig.MyURL = "http://" + appConfig.MyDomain
		} else {
			appConfig.MyURL = "https://" + appConfig.MyDomain
		}
	}
	g.Config = appConfig

	// check if it is test environment
	if os.Getenv("TESTENV") == "1" {
		g.TestEnv = true
	}
}

// InitServices initializes connections / sessions with different services
func InitServices(appDir string) {
	// get IP address of outbound interface
	g.MyIP = getOutboundIP(g.Config.PingIP)

	//
	// inititalize route cache
	//
	g.Routes = apiroutecache.InitSession(g.Config.RouteDb)

	//
	// inititalize mongodb for application use
	//
	g.Mgosesion = mango.InitSession(g.Config.Mongo)

	//
	// inititalize data cache
	//
	g.Cache = redicache.InitSession(g.Config.Redis)
	viewmodels.SetCacheSession(g.Config.Redis.ExpirationInMinute, g.Cache)

	//
	// initialize JWT Validator
	//
	//    Do not check tokens while unit testing
	if !g.TestEnv {
		g.JWTval = jwtauth.InitValidator(g.Config.JwtAuthConfig, appDir+"/public.pem")
	}

	//
	//initialize Api
	//
	//g.Api = apiclient.API{AllowInsecureSSL: !c.IsProduction, Timeout: 10}
	g.Api = apiclient.API{AllowInsecureSSL: false, Timeout: 10}

	// ----------------------
	// Initialize db connections
	//
	g.Db = pgsql.InitDbPool(g.Config.DBReader, g.Config.DBWriter)
	//g.Db.SetLogger(*g.Logger)

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

	switch g.Config.LogLevel {
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

	g.Logger = logger.New(logerror, loginfo, logdebug, true)

	// set output to file
	//create logs directory if not exist
	newpath := filepath.Join(appDir, "logs")
	logpath := filepath.Join(newpath, "app.log")
	os.Mkdir(newpath, os.ModePerm)
	var err error
	// check log file exist, if not then create
	if _, err = os.Stat(logpath); os.IsNotExist(err) {
		logFile, err = os.Create(logpath)
		if err != nil {
			g.Logger.Fatal().Err(err).Msg("log file create error")
		}
	} else {
		logFile, err = os.Open(logpath)
		if err != nil {
			g.Logger.Fatal().Err(err).Msg("log file open error")
		}
	}

	g.Logger.Output(logFile)
}

//AppCleanup cleanup existing connections and sessions
func AppCleanup() {
	g.Routes.Cleanup()
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
