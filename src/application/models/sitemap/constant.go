package sitemap // import "application/models/sitemap"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

const (
	// Resource type
	Resource = Type(`sitemap`)

	// ResourceIndex type
	ResourceIndex = Type(`sitemap-index`)
)

// Type of sitemap xml content
type Type string

// String convert Type to string
func (t Type) String() string { return string(t) }
