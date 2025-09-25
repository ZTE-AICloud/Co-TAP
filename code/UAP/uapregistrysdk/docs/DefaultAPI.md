# \DefaultAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RoutesGet**](DefaultAPI.md#RoutesGet) | **Get** /routes | Get all routes
[**RoutesPost**](DefaultAPI.md#RoutesPost) | **Post** /routes | Register a new route
[**RoutesRouteNameDelete**](DefaultAPI.md#RoutesRouteNameDelete) | **Delete** /routes/{routeName} | Delete a route
[**RoutesRouteNameGet**](DefaultAPI.md#RoutesRouteNameGet) | **Get** /routes/{routeName} | Get a route by name
[**RoutesRouteNamePut**](DefaultAPI.md#RoutesRouteNamePut) | **Put** /routes/{routeName} | Update a route
[**ServicesGet**](DefaultAPI.md#ServicesGet) | **Get** /services | Get all services
[**ServicesIdDelete**](DefaultAPI.md#ServicesIdDelete) | **Delete** /services/{id} | Delete a service
[**ServicesIdPatch**](DefaultAPI.md#ServicesIdPatch) | **Patch** /services/{id} | Update a service
[**ServicesIdRenewalPut**](DefaultAPI.md#ServicesIdRenewalPut) | **Put** /services/{id}/renewal | Renewal a service
[**ServicesPost**](DefaultAPI.md#ServicesPost) | **Post** /services | Create a new service
[**ServicesServiceNameGet**](DefaultAPI.md#ServicesServiceNameGet) | **Get** /services/{serviceName} | Get service by name



## RoutesGet

> []Route RoutesGet(ctx).Wait(wait).Index(index).Execute()

Get all routes

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	wait := "wait_example" // string | wait seconds (optional)
	index := "index_example" // string | index (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RoutesGet(context.Background()).Wait(wait).Index(index).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RoutesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RoutesGet`: []Route
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RoutesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRoutesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **wait** | **string** | wait seconds | 
 **index** | **string** | index | 

### Return type

[**[]Route**](Route.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RoutesPost

> Route RoutesPost(ctx).CreateRouteRequest(createRouteRequest).Execute()

Register a new route

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	createRouteRequest := *openapiclient.NewCreateRouteRequest("Name_example") // CreateRouteRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RoutesPost(context.Background()).CreateRouteRequest(createRouteRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RoutesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RoutesPost`: Route
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RoutesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRoutesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createRouteRequest** | [**CreateRouteRequest**](CreateRouteRequest.md) |  | 

### Return type

[**Route**](Route.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RoutesRouteNameDelete

> RoutesRouteNameDelete(ctx, routeName).Execute()

Delete a route

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	routeName := "routeName_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.RoutesRouteNameDelete(context.Background(), routeName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RoutesRouteNameDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**routeName** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRoutesRouteNameDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RoutesRouteNameGet

> Route RoutesRouteNameGet(ctx, routeName).Wait(wait).Index(index).Execute()

Get a route by name

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	routeName := "routeName_example" // string | 
	wait := "wait_example" // string | wait seconds (optional)
	index := "index_example" // string | index (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RoutesRouteNameGet(context.Background(), routeName).Wait(wait).Index(index).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RoutesRouteNameGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RoutesRouteNameGet`: Route
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RoutesRouteNameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**routeName** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRoutesRouteNameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **wait** | **string** | wait seconds | 
 **index** | **string** | index | 

### Return type

[**Route**](Route.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RoutesRouteNamePut

> Route RoutesRouteNamePut(ctx, routeName).UpdateRouteRequest(updateRouteRequest).Execute()

Update a route

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	routeName := "routeName_example" // string | 
	updateRouteRequest := *openapiclient.NewUpdateRouteRequest("Name_example") // UpdateRouteRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.RoutesRouteNamePut(context.Background(), routeName).UpdateRouteRequest(updateRouteRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.RoutesRouteNamePut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RoutesRouteNamePut`: Route
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.RoutesRouteNamePut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**routeName** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRoutesRouteNamePutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateRouteRequest** | [**UpdateRouteRequest**](UpdateRouteRequest.md) |  | 

### Return type

[**Route**](Route.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesGet

> []Service ServicesGet(ctx).Wait(wait).Index(index).Execute()

Get all services



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	wait := "wait_example" // string | wait seconds (optional)
	index := "index_example" // string | index (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ServicesGet(context.Background()).Wait(wait).Index(index).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ServicesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServicesGet`: []Service
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ServicesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiServicesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **wait** | **string** | wait seconds | 
 **index** | **string** | index | 

### Return type

[**[]Service**](Service.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesIdDelete

> ServicesIdDelete(ctx, id).Execute()

Delete a service



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	id := "id_example" // string | ID of the service to delete

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.ServicesIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ServicesIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the service to delete | 

### Other Parameters

Other parameters are passed through a pointer to a apiServicesIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesIdPatch

> Service ServicesIdPatch(ctx, id).PatchServiceRequest(patchServiceRequest).Execute()

Update a service



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	id := "id_example" // string | ID of the service to update
	patchServiceRequest := *openapiclient.NewPatchServiceRequest("Name_example", "Protocol_example", "Host_example", int32(123)) // PatchServiceRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ServicesIdPatch(context.Background(), id).PatchServiceRequest(patchServiceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ServicesIdPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServicesIdPatch`: Service
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ServicesIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the service to update | 

### Other Parameters

Other parameters are passed through a pointer to a apiServicesIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **patchServiceRequest** | [**PatchServiceRequest**](PatchServiceRequest.md) |  | 

### Return type

[**Service**](Service.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesIdRenewalPut

> ServicesIdRenewalPut(ctx, id).Execute()

Renewal a service



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	id := "id_example" // string | ID of the service to renewal

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.ServicesIdRenewalPut(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ServicesIdRenewalPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | ID of the service to renewal | 

### Other Parameters

Other parameters are passed through a pointer to a apiServicesIdRenewalPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesPost

> Service ServicesPost(ctx).CreateServiceRequest(createServiceRequest).Execute()

Create a new service



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	createServiceRequest := *openapiclient.NewCreateServiceRequest("Name_example", "Protocol_example", "Host_example", int32(123)) // CreateServiceRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ServicesPost(context.Background()).CreateServiceRequest(createServiceRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ServicesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServicesPost`: Service
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ServicesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiServicesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createServiceRequest** | [**CreateServiceRequest**](CreateServiceRequest.md) |  | 

### Return type

[**Service**](Service.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ServicesServiceNameGet

> []Service ServicesServiceNameGet(ctx, serviceName).Wait(wait).Index(index).Execute()

Get service by name



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ZTE-AICloud/Co-TAP/code/UAP/uapregistrysdk"
)

func main() {
	serviceName := "serviceName_example" // string | Name of the service to retrieve
	wait := "wait_example" // string | wait seconds (optional)
	index := "index_example" // string | index (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ServicesServiceNameGet(context.Background(), serviceName).Wait(wait).Index(index).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ServicesServiceNameGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ServicesServiceNameGet`: []Service
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ServicesServiceNameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**serviceName** | **string** | Name of the service to retrieve | 

### Other Parameters

Other parameters are passed through a pointer to a apiServicesServiceNameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **wait** | **string** | wait seconds | 
 **index** | **string** | index | 

### Return type

[**[]Service**](Service.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

