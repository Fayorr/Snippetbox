package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")
var ErrNoRecords = errors.New("models: no matching records found")