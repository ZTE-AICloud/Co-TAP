# CreateServiceRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Ephemeral** | Pointer to **bool** | 服务注册类型。true：临时注册；false：持久注册 | [optional] 
**Name** | **string** | 服务名称 | 
**Retries** | Pointer to **int32** | 连接失败重试次数 | [optional] 
**Protocol** | **string** | 传输协议类型 | 
**Host** | **string** | 后端服务域名/IP | 
**Port** | **int32** | 后端服务端口 | 
**Path** | Pointer to **string** | 服务路径 | [optional] 
**ConnectTimeout** | Pointer to **int32** | 连接超时时间（毫秒） | [optional] 
**WriteTimeout** | Pointer to **int32** | 写入超时时间（毫秒） | [optional] 
**ReadTimeout** | Pointer to **int32** | 读取超时时间（毫秒） | [optional] 
**Tags** | Pointer to **[]string** | 自定义标签（例如[namespacea,group&#x3D;x]） | [optional] 
**EphemeralCheck** | Pointer to [**EphemeralCheck**](EphemeralCheck.md) |  | [optional] 
**PersistentCheck** | Pointer to [**PersistentCheck**](PersistentCheck.md) |  | [optional] 
**AgentProtocol** | Pointer to **string** | Agent/Tool服务通信协议 | [optional] 
**AgentInfo** | Pointer to **map[string]interface{}** | 不同agent_protocol对应的特有内容 | [optional] 
**AgentInfoUrl** | Pointer to **string** | agent_info为空时，从agent_info_url获取agent_protocol内容 | [optional] 

## Methods

### NewCreateServiceRequest

`func NewCreateServiceRequest(name string, protocol string, host string, port int32, ) *CreateServiceRequest`

NewCreateServiceRequest instantiates a new CreateServiceRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateServiceRequestWithDefaults

`func NewCreateServiceRequestWithDefaults() *CreateServiceRequest`

NewCreateServiceRequestWithDefaults instantiates a new CreateServiceRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEphemeral

`func (o *CreateServiceRequest) GetEphemeral() bool`

GetEphemeral returns the Ephemeral field if non-nil, zero value otherwise.

### GetEphemeralOk

`func (o *CreateServiceRequest) GetEphemeralOk() (*bool, bool)`

GetEphemeralOk returns a tuple with the Ephemeral field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEphemeral

`func (o *CreateServiceRequest) SetEphemeral(v bool)`

SetEphemeral sets Ephemeral field to given value.

### HasEphemeral

`func (o *CreateServiceRequest) HasEphemeral() bool`

HasEphemeral returns a boolean if a field has been set.

### GetName

`func (o *CreateServiceRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateServiceRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateServiceRequest) SetName(v string)`

SetName sets Name field to given value.


### GetRetries

`func (o *CreateServiceRequest) GetRetries() int32`

GetRetries returns the Retries field if non-nil, zero value otherwise.

### GetRetriesOk

`func (o *CreateServiceRequest) GetRetriesOk() (*int32, bool)`

GetRetriesOk returns a tuple with the Retries field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetries

`func (o *CreateServiceRequest) SetRetries(v int32)`

SetRetries sets Retries field to given value.

### HasRetries

`func (o *CreateServiceRequest) HasRetries() bool`

HasRetries returns a boolean if a field has been set.

### GetProtocol

`func (o *CreateServiceRequest) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *CreateServiceRequest) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *CreateServiceRequest) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.


### GetHost

`func (o *CreateServiceRequest) GetHost() string`

GetHost returns the Host field if non-nil, zero value otherwise.

### GetHostOk

`func (o *CreateServiceRequest) GetHostOk() (*string, bool)`

GetHostOk returns a tuple with the Host field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHost

`func (o *CreateServiceRequest) SetHost(v string)`

SetHost sets Host field to given value.


### GetPort

`func (o *CreateServiceRequest) GetPort() int32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *CreateServiceRequest) GetPortOk() (*int32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *CreateServiceRequest) SetPort(v int32)`

SetPort sets Port field to given value.


### GetPath

`func (o *CreateServiceRequest) GetPath() string`

GetPath returns the Path field if non-nil, zero value otherwise.

### GetPathOk

`func (o *CreateServiceRequest) GetPathOk() (*string, bool)`

GetPathOk returns a tuple with the Path field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPath

`func (o *CreateServiceRequest) SetPath(v string)`

SetPath sets Path field to given value.

### HasPath

`func (o *CreateServiceRequest) HasPath() bool`

HasPath returns a boolean if a field has been set.

### GetConnectTimeout

`func (o *CreateServiceRequest) GetConnectTimeout() int32`

GetConnectTimeout returns the ConnectTimeout field if non-nil, zero value otherwise.

### GetConnectTimeoutOk

`func (o *CreateServiceRequest) GetConnectTimeoutOk() (*int32, bool)`

GetConnectTimeoutOk returns a tuple with the ConnectTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectTimeout

`func (o *CreateServiceRequest) SetConnectTimeout(v int32)`

SetConnectTimeout sets ConnectTimeout field to given value.

### HasConnectTimeout

`func (o *CreateServiceRequest) HasConnectTimeout() bool`

HasConnectTimeout returns a boolean if a field has been set.

### GetWriteTimeout

`func (o *CreateServiceRequest) GetWriteTimeout() int32`

GetWriteTimeout returns the WriteTimeout field if non-nil, zero value otherwise.

### GetWriteTimeoutOk

`func (o *CreateServiceRequest) GetWriteTimeoutOk() (*int32, bool)`

GetWriteTimeoutOk returns a tuple with the WriteTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWriteTimeout

`func (o *CreateServiceRequest) SetWriteTimeout(v int32)`

SetWriteTimeout sets WriteTimeout field to given value.

### HasWriteTimeout

`func (o *CreateServiceRequest) HasWriteTimeout() bool`

HasWriteTimeout returns a boolean if a field has been set.

### GetReadTimeout

`func (o *CreateServiceRequest) GetReadTimeout() int32`

GetReadTimeout returns the ReadTimeout field if non-nil, zero value otherwise.

### GetReadTimeoutOk

`func (o *CreateServiceRequest) GetReadTimeoutOk() (*int32, bool)`

GetReadTimeoutOk returns a tuple with the ReadTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadTimeout

`func (o *CreateServiceRequest) SetReadTimeout(v int32)`

SetReadTimeout sets ReadTimeout field to given value.

### HasReadTimeout

`func (o *CreateServiceRequest) HasReadTimeout() bool`

HasReadTimeout returns a boolean if a field has been set.

### GetTags

`func (o *CreateServiceRequest) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *CreateServiceRequest) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *CreateServiceRequest) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *CreateServiceRequest) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetEphemeralCheck

`func (o *CreateServiceRequest) GetEphemeralCheck() EphemeralCheck`

GetEphemeralCheck returns the EphemeralCheck field if non-nil, zero value otherwise.

### GetEphemeralCheckOk

`func (o *CreateServiceRequest) GetEphemeralCheckOk() (*EphemeralCheck, bool)`

GetEphemeralCheckOk returns a tuple with the EphemeralCheck field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEphemeralCheck

`func (o *CreateServiceRequest) SetEphemeralCheck(v EphemeralCheck)`

SetEphemeralCheck sets EphemeralCheck field to given value.

### HasEphemeralCheck

`func (o *CreateServiceRequest) HasEphemeralCheck() bool`

HasEphemeralCheck returns a boolean if a field has been set.

### GetPersistentCheck

`func (o *CreateServiceRequest) GetPersistentCheck() PersistentCheck`

GetPersistentCheck returns the PersistentCheck field if non-nil, zero value otherwise.

### GetPersistentCheckOk

`func (o *CreateServiceRequest) GetPersistentCheckOk() (*PersistentCheck, bool)`

GetPersistentCheckOk returns a tuple with the PersistentCheck field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPersistentCheck

`func (o *CreateServiceRequest) SetPersistentCheck(v PersistentCheck)`

SetPersistentCheck sets PersistentCheck field to given value.

### HasPersistentCheck

`func (o *CreateServiceRequest) HasPersistentCheck() bool`

HasPersistentCheck returns a boolean if a field has been set.

### GetAgentProtocol

`func (o *CreateServiceRequest) GetAgentProtocol() string`

GetAgentProtocol returns the AgentProtocol field if non-nil, zero value otherwise.

### GetAgentProtocolOk

`func (o *CreateServiceRequest) GetAgentProtocolOk() (*string, bool)`

GetAgentProtocolOk returns a tuple with the AgentProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentProtocol

`func (o *CreateServiceRequest) SetAgentProtocol(v string)`

SetAgentProtocol sets AgentProtocol field to given value.

### HasAgentProtocol

`func (o *CreateServiceRequest) HasAgentProtocol() bool`

HasAgentProtocol returns a boolean if a field has been set.

### GetAgentInfo

`func (o *CreateServiceRequest) GetAgentInfo() map[string]interface{}`

GetAgentInfo returns the AgentInfo field if non-nil, zero value otherwise.

### GetAgentInfoOk

`func (o *CreateServiceRequest) GetAgentInfoOk() (*map[string]interface{}, bool)`

GetAgentInfoOk returns a tuple with the AgentInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentInfo

`func (o *CreateServiceRequest) SetAgentInfo(v map[string]interface{})`

SetAgentInfo sets AgentInfo field to given value.

### HasAgentInfo

`func (o *CreateServiceRequest) HasAgentInfo() bool`

HasAgentInfo returns a boolean if a field has been set.

### GetAgentInfoUrl

`func (o *CreateServiceRequest) GetAgentInfoUrl() string`

GetAgentInfoUrl returns the AgentInfoUrl field if non-nil, zero value otherwise.

### GetAgentInfoUrlOk

`func (o *CreateServiceRequest) GetAgentInfoUrlOk() (*string, bool)`

GetAgentInfoUrlOk returns a tuple with the AgentInfoUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentInfoUrl

`func (o *CreateServiceRequest) SetAgentInfoUrl(v string)`

SetAgentInfoUrl sets AgentInfoUrl field to given value.

### HasAgentInfoUrl

`func (o *CreateServiceRequest) HasAgentInfoUrl() bool`

HasAgentInfoUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


