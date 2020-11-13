//go:generate go-enum -f=$GOFILE --marshal

package healthcheck

// HealthStatus defines the health statuses.
/* ENUM(
NotSet
OK
Warning
Critical
)
*/
type HealthStatus int
