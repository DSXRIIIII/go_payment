# \DefaultAPI

All URIs are relative to *https://127.0.0.1/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CustomerCustomerIDOrdersOrderIDGet**](DefaultAPI.md#CustomerCustomerIDOrdersOrderIDGet) | **Get** /customer/{customerID}/orders/{orderID} | 
[**CustomerCustomerIDOrdersPost**](DefaultAPI.md#CustomerCustomerIDOrdersPost) | **Post** /customer/{customerID}/orders | 



## CustomerCustomerIDOrdersOrderIDGet

> Order CustomerCustomerIDOrdersOrderIDGet(ctx, customerID, orderID).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	customerID := "customerID_example" // string | 
	orderID := "orderID_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.CustomerCustomerIDOrdersOrderIDGet(context.Background(), customerID, orderID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CustomerCustomerIDOrdersOrderIDGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CustomerCustomerIDOrdersOrderIDGet`: Order
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CustomerCustomerIDOrdersOrderIDGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**customerID** | **string** |  | 
**orderID** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCustomerCustomerIDOrdersOrderIDGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Order**](Order.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CustomerCustomerIDOrdersPost

> Order CustomerCustomerIDOrdersPost(ctx, customerID).CreateOrderRequest(createOrderRequest).Execute()





### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	customerID := "customerID_example" // string | 
	createOrderRequest := *openapiclient.NewCreateOrderRequest("CustomerID_example", []openapiclient.ItemWithQuantity{*openapiclient.NewItemWithQuantity()}) // CreateOrderRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.CustomerCustomerIDOrdersPost(context.Background(), customerID).CreateOrderRequest(createOrderRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.CustomerCustomerIDOrdersPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CustomerCustomerIDOrdersPost`: Order
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.CustomerCustomerIDOrdersPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**customerID** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCustomerCustomerIDOrdersPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createOrderRequest** | [**CreateOrderRequest**](CreateOrderRequest.md) |  | 

### Return type

[**Order**](Order.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

