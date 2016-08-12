package app

import "github.com/revel/revel"

// AppConfig is a global variable defining the directory to be processed by the
// app
var conf AppConfig

var Context *WebManager

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		CORSFilter,
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.CompressFilter,    // Compress the result.
		revel.ActionInvoker,     // Invoke the action.
	}

	revel.OnAppStart(InitAppConfig)
	revel.OnAppStart(InitApp)
}

var CORSFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")
	c.Response.Out.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Response.Out.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Stop here for a Preflighted OPTIONS request.
	if c.Request.Method == "OPTIONS" {
		return
	}

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
