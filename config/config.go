package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/samtech09/api-template/psql"
	"github.com/samtech09/apiroutecache"
	"github.com/samtech09/jwtauth"
	"github.com/samtech09/redicache"
)

// const ImageAPIBaseURL = "http://192.168.60.206/"
// const AuthAPIBaseURL = "http://74.124.24.16/mibsws/"

//appConfig hold paramater to run server
type appConfig struct {
	//HttpPort to listen for requests
	ListenPort int

	//PingIP is ipaddress and port that is accessible from the
	//  machine where app is running
	//  it is used to get internal IP address of current system
	PingIP string

	//DisableGzipset flag to disable gzip compression on response
	// bydefault it is enabled
	DisableGzip bool

	//EnableSSL - enables SSL/TLS
	//  only applicable to Development environment
	//  production is forced to use SSL
	DisableSSL bool
	//SSLCertFile is path to SSL certificate file
	SSLCertFile string
	//SSLCertFile is path to SSL certificate KEY file
	SSLKeyFile string

	//database config for Read-only database
	DBReader psql.PDbConfig

	//database config for Read-write database
	DBWriter psql.PDbConfig

	//MibsAPIBaseURL is base URL of API to get MIBS related data
	MibsAPIBaseURL string

	//Authentication serverconfig
	//  Only used for data-sync (refresh routes)
	AuthSrv Authserver

	//MongoDB config
	Mongo apiroutecache.MongoConfig

	//redis cache
	Redis redicache.RedisConfig

	//JWT token verification
	JwtAuthConfig jwtauth.ValidatorConfig

	//server own domain name
	MyDomain string

	//LogLevel tells which data should be logged into file
	/* Possible values are
	* fatal
	* errors = errors + fatal
	* info = errors + info
	* debug = debug + info
	*
	*   default is info
	 */
	LogLevel string
}

//Authserver config
type Authserver struct {
	//AuthServerURL is Authentication server which provides JWT tokens
	AuthServerURL            string
	AuthTokenEndPoint        string
	AuthRefreshTokenEndPoint string
	//AuthRevokeEndPoint is endpoint to revoke refresh-token
	AuthRevokeEndPoint string
	//AuthRouteAccessControlEndPoint provides list of routes for given scopes
	AuthRouteAccessControlEndPoint string
	AuthServerVendorID             string
}

// //RedisCache is only used cache LabST items
// type RedisCache struct {
// 	//RedisServerHost is name of machine to connect for Redis Cahce
// 	Host string
// 	//RedisServerPort is port of machine to connect for Redis Cahce
// 	Port               int
// 	DB                 int // from 0 to 16
// 	Pwd                string
// 	KeyPrefix          string
// 	ExpirationInMinute int
// }

var (
	//AppConfig is application/server configuration
	AppConfig appConfig
	//IsProduction is flag to run server in production mode
	IsProduction bool

	//MyURL is this sever's own URL e.g. https://mshopapi.mahendras.org.org etc
	MyURL string
	//OstImageFolder = "ostimgs"

	//ImgCDNPlaceholder is used to create actual Image url by replacing placeholder with URL prefix
	// as it is being saved alsong with Questiondata, passage, instructions
	// so *** DO NOT CHANGE IT ***
	ImgCDNPlaceholder = "[imgcdn]"
)

//Initconfig read config file and initialize AppConfig struct
func Initconfig(prod *bool, confFolderPath string) {
	//set if productin mode is true
	IsProduction = *prod

	conffile := confFolderPath + "conf.dev.json"
	if IsProduction {
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
	AppConfig = appConfig{}
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatal("Cannot parse config.\n", err)
	}

	if IsProduction {
		MyURL = "https://" + AppConfig.MyDomain
	} else {
		if AppConfig.DisableSSL {
			MyURL = "http://" + AppConfig.MyDomain
		} else {
			MyURL = "https://" + AppConfig.MyDomain
		}
	}

	// if !IsProduction {
	// 	fmt.Printf("Config: %v\n", AppConfig)
	// }
}
