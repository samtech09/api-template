package global

import (
	"net"

	"github.com/rs/zerolog"
	"github.com/samtech09/api-template/mango"
	"github.com/samtech09/api-template/psql"
	"github.com/samtech09/apiclient"
	"github.com/samtech09/apiroutecache"
	"github.com/samtech09/jwtauth"
	"github.com/samtech09/redicache"
)

var (
	Api       apiclient.API
	Logger    *zerolog.Logger
	Db        *psql.Db
	Routes    *apiroutecache.MongoSession
	Mgosesion *mango.MongoSession
	JWTval    *jwtauth.Validator
	Cache     *redicache.RedisSession
	Config    appConfig
	//IsProduction is flag to run server in production mode
	IsProduction bool

	//TestEnv tell if running in Test mode for unit testing
	//Set environment variable TESTENV=1 to set it
	TestEnv bool

	//MyURL is this sever's own URL e.g. https://mshopapi.mahendras.org.org etc
	MyURL string
	MyIP  net.IP
)

//NewConfig creates new blank config
func NewConfig() appConfig {
	return appConfig{}
}

//AppConfig hold paramater to run server
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

	//ExtApiBaseUrl is base URL of external API (if any required)
	ExtApiBaseUrl string

	//Authentication serverconfig
	//  Only used for data-sync (refresh routes)
	AuthSrv authserver

	//RouteDb is API route cache in mongodb
	RouteDb apiroutecache.MongoConfig

	//Mongo is config for MongoDb
	Mongo mango.MongoConfig

	//redis cache
	Redis redicache.RedisConfig

	//JWT token verification
	JwtAuthConfig jwtauth.ValidatorConfig

	//server's domain name
	MyDomain string
	MyURL    string

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
type authserver struct {
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
