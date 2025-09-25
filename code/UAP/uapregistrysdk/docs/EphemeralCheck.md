# EphemeralCheck

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CheckType** | **string** | 健康检查类型：TTL | [default to "TTL"]
**RenewalDeleteTimeout** | Pointer to **string** | 服务超期未更新则删除 | [optional] [default to "60s"]
**RenewalInterval** | Pointer to **string** | 服务更新周期 | [optional] [default to "15s"]
**RenewalUnhealthyTimeout** | Pointer to **string** | 服务超期未更新则不健康 | [optional] [default to "30s"]

## Methods

### NewEphemeralCheck

`func NewEphemeralCheck(checkType string, ) *EphemeralCheck`

NewEphemeralCheck instantiates a new EphemeralCheck object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEphemeralCheckWithDefaults

`func NewEphemeralCheckWithDefaults() *EphemeralCheck`

NewEphemeralCheckWithDefaults instantiates a new EphemeralCheck object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCheckType

`func (o *EphemeralCheck) GetCheckType() string`

GetCheckType returns the CheckType field if non-nil, zero value otherwise.

### GetCheckTypeOk

`func (o *EphemeralCheck) GetCheckTypeOk() (*string, bool)`

GetCheckTypeOk returns a tuple with the CheckType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCheckType

`func (o *EphemeralCheck) SetCheckType(v string)`

SetCheckType sets CheckType field to given value.


### GetRenewalDeleteTimeout

`func (o *EphemeralCheck) GetRenewalDeleteTimeout() string`

GetRenewalDeleteTimeout returns the RenewalDeleteTimeout field if non-nil, zero value otherwise.

### GetRenewalDeleteTimeoutOk

`func (o *EphemeralCheck) GetRenewalDeleteTimeoutOk() (*string, bool)`

GetRenewalDeleteTimeoutOk returns a tuple with the RenewalDeleteTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRenewalDeleteTimeout

`func (o *EphemeralCheck) SetRenewalDeleteTimeout(v string)`

SetRenewalDeleteTimeout sets RenewalDeleteTimeout field to given value.

### HasRenewalDeleteTimeout

`func (o *EphemeralCheck) HasRenewalDeleteTimeout() bool`

HasRenewalDeleteTimeout returns a boolean if a field has been set.

### GetRenewalInterval

`func (o *EphemeralCheck) GetRenewalInterval() string`

GetRenewalInterval returns the RenewalInterval field if non-nil, zero value otherwise.

### GetRenewalIntervalOk

`func (o *EphemeralCheck) GetRenewalIntervalOk() (*string, bool)`

GetRenewalIntervalOk returns a tuple with the RenewalInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRenewalInterval

`func (o *EphemeralCheck) SetRenewalInterval(v string)`

SetRenewalInterval sets RenewalInterval field to given value.

### HasRenewalInterval

`func (o *EphemeralCheck) HasRenewalInterval() bool`

HasRenewalInterval returns a boolean if a field has been set.

### GetRenewalUnhealthyTimeout

`func (o *EphemeralCheck) GetRenewalUnhealthyTimeout() string`

GetRenewalUnhealthyTimeout returns the RenewalUnhealthyTimeout field if non-nil, zero value otherwise.

### GetRenewalUnhealthyTimeoutOk

`func (o *EphemeralCheck) GetRenewalUnhealthyTimeoutOk() (*string, bool)`

GetRenewalUnhealthyTimeoutOk returns a tuple with the RenewalUnhealthyTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRenewalUnhealthyTimeout

`func (o *EphemeralCheck) SetRenewalUnhealthyTimeout(v string)`

SetRenewalUnhealthyTimeout sets RenewalUnhealthyTimeout field to given value.

### HasRenewalUnhealthyTimeout

`func (o *EphemeralCheck) HasRenewalUnhealthyTimeout() bool`

HasRenewalUnhealthyTimeout returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


