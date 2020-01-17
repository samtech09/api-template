package viewmodels

import (
	"time"

	"github.com/samtech09/redicache"
)

var _defExpirationInMinutes int
var _cache *redicache.RedisSession

//SetCacheSession set session object where cache items could be added
func SetCacheSession(expinminute int, c *redicache.RedisSession) {
	_defExpirationInMinutes = expinminute
	_cache = c
}

//RegisterModelsForCaching register models to be cached in Redis Cache
func RegisterModelsForCaching() {
	_cache.RegisterCandidate(DbUser{}, "Sample cache candidate to show usage")

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
	return time.Minute * time.Duration(_defExpirationInMinutes)
}
