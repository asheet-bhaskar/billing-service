##  Billing Service
Go application that manages the billing lifecycle and workflows. 

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


### Endpoints
#### create customer 
```
curl -X POST 'localhost:4000/customers' -d '{"FirstName":"","LastName":"","Email":""}'
```

#### get customer by id
```
curl -X GET 'localhost:4000/customers/:id'
```

#### create currency
```
curl  -X POST 'localhost:4000/currencies' -d '{"Code":"","Name":"","Symbol":""}'
```

#### get currency by id
```
curl -X GET  'localhost:4000/currencies/:id'
```

#### create bill
```
curl -X POST 'localhost:4000/bills' -d '{"Description":"","CustomerID":"","CurrencyCode":"","PeriodStart":"2009-11-10T23:00:00Z","PeriodEnd":"2009-11-10T23:00:00Z"}'
```

#### get bill by id
```
curl -X GET 'localhost:4000/bills/:id'
```

#### add line item to bill
```
curl -X POST 'localhost:4000/bills/items' -d '{"BillID":"","Description":"","Amount":0}'
```

#### remove line item from bill
```
curl -X PUT 'localhost:4000/bills/:billID/items/:itemID'
```

#### close bill by id
```
curl -X PUT 'localhost:4000/bills/:id/close'
```

#### get invoice
```
curl -X GET 'localhost:4000/bills/:id/invoice'
```

### Further Improvements
* Create distinction among request and corresponding models, like bill and billRequest
* Make application dependencies configurable
  * temporal server
  * postgresql
  * etc
* Add integration tests

