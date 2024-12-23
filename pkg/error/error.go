package error

import "errors"

var BillNotFoundError = errors.New("Bill not found")
var BillClosedError = errors.New("Bill is closed")
var LineItemNotFoundError = errors.New("Line item not found")
var CustomerNotFoundError = errors.New("Customer not found")
var CurrencyNotFoundError = errors.New("Currency not found")

var BillAlreadyExistError = errors.New("Bill already exist")
var LineItemAlreadyExistError = errors.New("Line item already exist")
var CustomerAlreadyExistError = errors.New("Customer already exist")
var CurrencyAlreadyExistError = errors.New("Currency already exist")
