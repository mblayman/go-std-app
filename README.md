# go-std-app

Make a Go app using nothing but the standard library

## Run

```
go run app.go
# Stop with Ctrl+C
```

Things to explore:

* Data
    * no native db option
    * is there a db driver for a CSV file?
    * What's up with BoltDB?
    * If we allowed one driver, what's it like to work with SQLite?
* middleware
* auth - r.Context to get a value
* Handle response compression? What does the default server cover?
