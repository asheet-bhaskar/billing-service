##  Billing Service


### Dependencies 
* [golang](https://go.dev/doc/install), version 1.22.5
* [gorm](https://gorm.io/)
* [encore](https://encore.dev/docs/ts/install), version v1.45.1
* [temporal](https://learn.temporal.io/getting_started/go/dev_environment/), version 1.1.2
* [postgresql]()

### Develoment Setup
* `git clone https://github.com/asheet-bhaskar/billing-service`
* `go mod tidy`
* start temporal server 
  * `temporal server start-dev`
* start encore server
  * `cd billing-service`
  * `encore run`
* open encore and temporal dashboards in web browser
  * [encore dashboard](http://localhost:9400/)
  * [temporal dashboard](http://localhost:8233/)
* Run tests 
  * `go test ./...`
* Database shell
  * `encore db shell billing_service` # Development DB
  * `encore db shell billing_service_test` # Test DB

### TODOs
* Make application dependencies configurable
  * temporal server
  * postgresql
  * etc
* Add integration tests

