# cfmgo
A MongoDB integration package for Cloud Foundry

## Overview

`cfmgo` is a package to assist you in connecting Go applications running on Cloud Foundry to MongoDB.  

## Usage

`go get github.com/pivotal-pez/cfmgo`

```go
appEnv, _ := cfenv,Current() //relies on github.com/cloudfoundry-community/go-cfenv
serviceName := os.Getenv("DB_NAME")
serviceURIName := os.Getenv("DB_URI")
serviceURI := cfmgo.GetServiceBinding(serviceName, serviceURIName, appEnv)
collection := cfmgo.Connect(cfmgo.NewCollectionDialer, serviceURI, "my-collection")
```

### cfmgo/params

`params` is a package that extracts query parameters from a request to be used in cfmgo.Collection.Find() operations.

#### Example
```go
func ListInventoryItemsHandler(collection cfmgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()

		params := params.Extract(req.URL.Query())

		items := make([]RedactedInventoryItem, 0)

		if count, err := collection.Find(params, &items); err == nil {
			Formatter().JSON(w, http.StatusOK, wrapper.Collection(&items, count))
		} else {
			Formatter().JSON(w, http.StatusNotFound, wrapper.Error(err.Error()))
		}
	}
}
```

### cfmgo/wrapper

`wrapper` is a simple helper to wrap API response data and errors in a consistent structure.  

```go
//ResponseWrapper provides a standard structure for API responses.
type ResponseWrapper struct {
	//Status indicates the result of a request as "success" or "error"
	Status string `json:"status"`
	//Data holds the payload of the response
	Data interface{} `json:"data,omitempty"`
	//Message contains the nature of an error
	Message string `json:"message,omitempty"`
	//Count contains the number of records in the result set
	Count int `json:"count,omitempty"`
}
```
#### Examples

`wrapper.Error(err)` yields:

```json
{
"status": "error",
"message": "error message text"
}
```

`wrapper.One(&someRecord)` yields:

```json
{
"status": "success",
"data": {
	"id": 1,
	"name": "fluffy"
	}
}
```

`wrapper.Collection(&someResults, count)` yields:

```json
{
"status": "success",
"data": [
	{
	"id": 1,
	"name": "fluffy"
	},
	{
	"id": 2,
	"name": "thiggy"
	}
	],
"count": 2
}
```

Note: `count` represents the total number of matching records from the query, not the number of records returned in the result set.  That number is governed by `limit`.
