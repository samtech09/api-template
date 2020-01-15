package global

import (
	"net"

	"github.com/rs/zerolog"
	"github.com/samtech09/api-template/psql"
	"github.com/samtech09/apiclient"
	"github.com/samtech09/apiroutecache"
	"github.com/samtech09/jwtauth"
	"github.com/samtech09/redicache"
)

var (
	Api       apiclient.API
	Logger    *zerolog.Logger
	MyIP      net.IP
	Db        *psql.Db
	Mgosesion *apiroutecache.MongoSession
	JWTval    *jwtauth.Validator
	Cache     *redicache.RedisSession
	TestEnv   bool
)
