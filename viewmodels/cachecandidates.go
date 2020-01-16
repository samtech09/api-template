package viewmodels

import (
	"time"

	g "github.com/samtech09/api-template/global"
)

//RegisterModelsForCaching register models to be cached in Redis Cache
func RegisterModelsForCaching() {
	g.Cache.RegisterCandidate(DbUser{}, "Sample cache candidate to show usage")

	// register more candidates ...
}

/*
 *
 * Implement Interface Methods to satisfy as valid CacheCandidate
 *
 */

//GetKey - parentid1 should be publishedfor::int
func (s DbUser) GetKey(parentid1, parentid2 string) string {
	return "USER:" + parentid1 + ":" + parentid2
}

//GetMasterKey - returns patter for key, just for documentation
func (s DbUser) GetMasterKey() string {
	return "USER:<parent1>:<parent2>"
}

//GetExpiration returns expiration duration of key
func (s DbUser) GetExpiration() time.Duration {
	return time.Minute * time.Duration(g.Config.Redis.ExpirationInMinute)
}
