package layouts

import "github.com/alessandrofascini/log4go/pkg"

// externalLayoutPool this variable is used to store custom layouts.
// It is a map: a key represents the name of the layout and its value represents the function used to create the layout.
var externalLayoutPool = make(map[string]*func(*pkg.LayoutConfig) pkg.Layout)

// AddExternalLayout add a custom Layout function
func AddExternalLayout(name string, customLayout *func(*pkg.LayoutConfig) pkg.Layout) {
	externalLayoutPool[name] = customLayout
}

// ClearLayoutPool empty the pool
func ClearLayoutPool() {
	externalLayoutPool = make(map[string]*func(*pkg.LayoutConfig) pkg.Layout)
}
