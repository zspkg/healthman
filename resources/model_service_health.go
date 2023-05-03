/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ServiceHealth struct {
	Key
	Attributes ServiceHealthAttributes `json:"attributes"`
}
type ServiceHealthResponse struct {
	Data     ServiceHealth `json:"data"`
	Included Included      `json:"included"`
}

type ServiceHealthListResponse struct {
	Data     []ServiceHealth `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustServiceHealth - returns ServiceHealth from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustServiceHealth(key Key) *ServiceHealth {
	var serviceHealth ServiceHealth
	if c.tryFindEntry(key, &serviceHealth) {
		return &serviceHealth
	}
	return nil
}
