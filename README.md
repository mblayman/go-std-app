# go-std-app

Make a Go app using nothing but the standard library

Things to explore:

* Data
    * no native db option
    * What's up with BoltDB?
    * If we allowed one driver, what's it like to work with SQLite?
* middleware
* auth - r.Context to get a value
* how are static files served? Returning binary files directly.
* cookies?
* HandlerFunc that calls another HTTP endpoint with an HTTP client
* Handle response compression? What does the default server cover?
