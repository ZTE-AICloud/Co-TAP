# PersistentCheck

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CheckType** | **string** | 健康检查类型：HTTP/HTTPS | 
**CheckInterval** | Pointer to **string** | 健康检查时间间隔 | [optional] [default to "15s"]
**CheckTimeout** | Pointer to **string** | 健康检查超时时间 | [optional] [default to "5s"]
**CheckUnhealthyTimeout** | Pointer to **string** | 服务超期健康检查一直未成功则不健康 | [optional] [default to "60s"]
**CheckHttpUrl** | **string** | 健康检查URL地址，示例：/healthcheck | 
**CheckHttpMethod** | Pointer to **string** | HTTP, HTTPS健康检查使用的HTTP 方法：GET, HEAD, OPTIONS | [optional] 

## Methods

### NewPersistentCheck

`func NewPersistentCheck(checkType string, checkHttpUrl string, ) *PersistentCheck`

NewPersistentCheck instantiates a new PersistentCheck object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPersistentCheckWithDefaults

`func NewPersistentCheckWithDefaults() *PersistentCheck`

NewPersistentCheckWithDefaults instantiates a new PersistentCheck object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCheckType

`func (o *PersistentCheck) GetCheckType() string`

GetCheckType returns the CheckType field if non-nil, zero value otherwise.

### GetCheckTypeOk

`func (o *PersistentCheck) GetCheckTypeOk() (*string, bool)`

GetCheckTypeOk returns a tuple with the CheckType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCheckType

`func (o *PersistentCheck) SetCheckType(v string)`

SetCheckType sets CheckType field to given value.


### GetCheckInterval

`func (o *PersistentCheck) GetCheckInterval() string`

GetCheckInterval returns the CheckInterval field if non-nil, zero value otherwise.

### GetCheckIntervalOk

`func (o *PersistentCheck) GetCheckIntervalOk() (*string, bool)`

GetCheckIntervalOk returns a tuple with the CheckInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCheckInterval

`func (o *PersistentCheck) SetCheckInterval(v string)`

SetCheckInterval sets CheckInterval field to given value.

### HasCheckInterval

`func (o *PersistentCheck) HasCheckInterval() bool`

HasCheckInterval returns a boolean if a field has been set.

### GetCheckTimeout

`func (o *PersistentCheck) GetCheckTimeout() string`

GetCheckTimeout returns the CheckTimeout field if non-nil, zero value otherwise.

### GetCheckTimeoutOk

`func (o *PersistentCheck) GetCheckTimeoutOk() (*string, bool)`

GetCheckTimeoutOk returns a tuple with the CheckTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCheckTimeout

`func (o *PersistentCheck) SetCheckTimeout(v string)`

SetCheckTimeout sets CheckTimeout field to given value.

### HasCheckTimeout

`func (o *PersistentCheck) HasCheckTimeout() bool`

HasCheckTimeout returns a boolean if a field has been set.

### GetCheckUnhealthyTimeout

`func (o *PersistentCheck) GetCheckUnhealthyTimeout() string`

GetCheckUnhealthyTimeout returns the CheckUnhealthyTimeout field if non-nil, zero value otherwise.

### GetCheckUnhealthyTimeoutOk

`func (o *PersistentCheck) GetCheckUnhealthyTimeoutOk() (*string, bool)`

GetCheckUnhealthyTimeoutOk returns a tuple with the CheckUnhealthyTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCheckUnhealthyTimeout

`func (o *PersistentCheck) SetCheckUnhealthyTimeout(v string)`

SetCheckUnhealthyTimeout sets CheckUnhealthyTimeout field to given value.

### HasCheckUnhealthyTimeout

`func (o *PersistentCheck) HasCheckUnhealthyTimeout() bool`

HasCheckUnhealthyTimeout returns a boolean if a field has been set.

### GetCheckHttpUrl

`func (o *PersistentCheck) GetCheckHttpUrl() string`

GetCheckHttpUrl returns the CheckHttpUrl field if non-nil, zero value otherwise.

### GetCheckHttpUrlOk

`func (o *PersistentCheck) GetCheckHttpUrlOk() (*string, bool)`

GetCheckHttpUrlOk returns a tuple with the CheckHttpUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCheckHttpUrl

`func (o *PersistentCheck) SetCheckHttpUrl(v string)`

SetCheckHttpUrl sets CheckHttpUrl field to given value.


### GetCheckHttpMethod

`func (o *PersistentCheck) GetCheckHttpMethod() string`

GetCheckHttpMethod returns the CheckHttpMethod field if non-nil, zero value otherwise.

### GetCheckHttpMethodOk

`func (o *PersistentCheck) GetCheckHttpMethodOk() (*string, bool)`

GetCheckHttpMethodOk returns a tuple with the CheckHttpMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCheckHttpMethod

`func (o *PersistentCheck) SetCheckHttpMethod(v string)`

SetCheckHttpMethod sets CheckHttpMethod field to given value.

### HasCheckHttpMethod

`func (o *PersistentCheck) HasCheckHttpMethod() bool`

HasCheckHttpMethod returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


