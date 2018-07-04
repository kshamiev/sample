package options // import "application/modules/rendering/options"

// Option Templating settings
type Option struct {
	// Directory to load templates. Default is "assets/www"
	Directory string

	// Reload to reload templates everytime. Default is false - сaching templates in memory
	Reload bool
}
