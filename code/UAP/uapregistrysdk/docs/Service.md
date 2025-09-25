# Service

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Ephemeral** | Pointer to **bool** | 服务注册类型。true：临时注册；false：持久注册 | [optional] 
**Id** | Pointer to **string** | 服务唯一标识符 | [optional] 
**CreatedAt** | Pointer to **int64** | 记录创建时间（由注册中心设置） | [optional] 
**UpdatedAt** | Pointer to **int64** | 记录最后更新时间（由注册中心设置） | [optional] 
**Index** | Pointer to **int64** | 服务更新序号 | [optional] 
**Name** | **string** | 服务名称 | 
**Retries** | Pointer to **int32** | 连接失败重试次数 | [optional] 
**Protocol** | **string** | 传输协议类型 | 
**Host** | **string** | 后端服务域名/IP | 
**Port** | **int32** | 后端服务端口 | 
**Path** | Pointer to **string** | 服务路径 | [optional] 
**ConnectTimeout** | Pointer to **int32** | 连接超时时间（毫秒） | [optional] 
**WriteTimeout** | Pointer to **int32** | 写入超时时间（毫秒） | [optional] 
**ReadTimeout** | Pointer to **int32** | 读取超时时间（毫秒） | [optional] 
**Tags** | Pointer to **[]string** | 自定义标签（例如[namespace&#x3D;a,group&#x3D;x]） | [optional] 
**EphemeralCheck** | Pointer to [**EphemeralCheck**](EphemeralCheck.md) |  | [optional] 
**PersistentCheck** | Pointer to [**PersistentCheck**](PersistentCheck.md) |  | [optional] 
**AgentProtocol** | Pointer to **string** | Agent/Tool服务通信协议 | [optional] 
**AgentInfo** | Pointer to **map[string]interface{}** | 不同agent_protocol对应的特有内容 | [optional] 
**AgentInfoUrl** | Pointer to **string** | agent_info为空时，从agent_info_url获取agent_protocol内容 | [optional] 

## Methods

### NewService

`func NewService(name string, protocol string, host string, port int32, ) *Service`

NewService instantiates a new Service object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServiceWithDefaults

`func NewServiceWithDefaults() *Service`

NewServiceWithDefaults instantiates a new Service object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEphemeral

`func (o *Service) GetEphemeral() bool`

GetEphemeral returns the Ephemeral field if non-nil, zero value otherwise.

### GetEphemeralOk

`func (o *Service) GetEphemeralOk() (*bool, bool)`

GetEphemeralOk returns a tuple with the Ephemeral field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEphemeral

`func (o *Service) SetEphemeral(v bool)`

SetEphemeral sets Ephemeral field to given value.

### HasEphemeral

`func (o *Service) HasEphemeral() bool`

HasEphemeral returns a boolean if a field has been set.

### GetId

`func (o *Service) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Service) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Service) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Service) HasId() bool`

HasId returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Service) GetCreatedAt() int64`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Service) GetCreatedAtOk() (*int64, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Service) SetCreatedAt(v int64)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Service) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *Service) GetUpdatedAt() int64`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Service) GetUpdatedAtOk() (*int64, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Service) SetUpdatedAt(v int64)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *Service) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetIndex

`func (o *Service) GetIndex() int64`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *Service) GetIndexOk() (*int64, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *Service) SetIndex(v int64)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *Service) HasIndex() bool`

HasIndex returns a boolean if a field has been set.

### GetName

`func (o *Service) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Service) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Service) SetName(v string)`

SetName sets Name field to given value.


### GetRetries

`func (o *Service) GetRetries() int32`

GetRetries returns the Retries field if non-nil, zero value otherwise.

### GetRetriesOk

`func (o *Service) GetRetriesOk() (*int32, bool)`

GetRetriesOk returns a tuple with the Retries field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetries

`func (o *Service) SetRetries(v int32)`

SetRetries sets Retries field to given value.

### HasRetries

`func (o *Service) HasRetries() bool`

HasRetries returns a boolean if a field has been set.

### GetProtocol

`func (o *Service) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *Service) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *Service) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.


### GetHost

`func (o *Service) GetHost() string`

GetHost returns the Host field if non-nil, zero value otherwise.

### GetHostOk

`func (o *Service) GetHostOk() (*string, bool)`

GetHostOk returns a tuple with the Host field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHost

`func (o *Service) SetHost(v string)`

SetHost sets Host field to given value.


### GetPort

`func (o *Service) GetPort() int32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *Service) GetPortOk() (*int32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *Service) SetPort(v int32)`

SetPort sets Port field to given value.


### GetPath

`func (o *Service) GetPath() string`

GetPath returns the Path field if non-nil, zero value otherwise.

### GetPathOk

`func (o *Service) GetPathOk() (*string, bool)`

GetPathOk returns a tuple with the Path field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPath

`func (o *Service) SetPath(v string)`

SetPath sets Path field to given value.

### HasPath

`func (o *Service) HasPath() bool`

HasPath returns a boolean if a field has been set.

### GetConnectTimeout

`func (o *Service) GetConnectTimeout() int32`

GetConnectTimeout returns the ConnectTimeout field if non-nil, zero value otherwise.

### GetConnectTimeoutOk

`func (o *Service) GetConnectTimeoutOk() (*int32, bool)`

GetConnectTimeoutOk returns a tuple with the ConnectTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectTimeout

`func (o *Service) SetConnectTimeout(v int32)`

SetConnectTimeout sets ConnectTimeout field to given value.

### HasConnectTimeout

`func (o *Service) HasConnectTimeout() bool`

HasConnectTimeout returns a boolean if a field has been set.

### GetWriteTimeout

`func (o *Service) GetWriteTimeout() int32`

GetWriteTimeout returns the WriteTimeout field if non-nil, zero value otherwise.

### GetWriteTimeoutOk

`func (o *Service) GetWriteTimeoutOk() (*int32, bool)`

GetWriteTimeoutOk returns a tuple with the WriteTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWriteTimeout

`func (o *Service) SetWriteTimeout(v int32)`

SetWriteTimeout sets WriteTimeout field to given value.

### HasWriteTimeout

`func (o *Service) HasWriteTimeout() bool`

HasWriteTimeout returns a boolean if a field has been set.

### GetReadTimeout

`func (o *Service) GetReadTimeout() int32`

GetReadTimeout returns the ReadTimeout field if non-nil, zero value otherwise.

### GetReadTimeoutOk

`func (o *Service) GetReadTimeoutOk() (*int32, bool)`

GetReadTimeoutOk returns a tuple with the ReadTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadTimeout

`func (o *Service) SetReadTimeout(v int32)`

SetReadTimeout sets ReadTimeout field to given value.

### HasReadTimeout

`func (o *Service) HasReadTimeout() bool`

HasReadTimeout returns a boolean if a field has been set.

### GetTags

`func (o *Service) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *Service) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *Service) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *Service) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetEphemeralCheck

`func (o *Service) GetEphemeralCheck() EphemeralCheck`

GetEphemeralCheck returns the EphemeralCheck field if non-nil, zero value otherwise.

### GetEphemeralCheckOk

`func (o *Service) GetEphemeralCheckOk() (*EphemeralCheck, bool)`

GetEphemeralCheckOk returns a tuple with the EphemeralCheck field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEphemeralCheck

`func (o *Service) SetEphemeralCheck(v EphemeralCheck)`

SetEphemeralCheck sets EphemeralCheck field to given value.

### HasEphemeralCheck

`func (o *Service) HasEphemeralCheck() bool`

HasEphemeralCheck returns a boolean if a field has been set.

### GetPersistentCheck

`func (o *Service) GetPersistentCheck() PersistentCheck`

GetPersistentCheck returns the PersistentCheck field if non-nil, zero value otherwise.

### GetPersistentCheckOk

`func (o *Service) GetPersistentCheckOk() (*PersistentCheck, bool)`

GetPersistentCheckOk returns a tuple with the PersistentCheck field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPersistentCheck

`func (o *Service) SetPersistentCheck(v PersistentCheck)`

SetPersistentCheck sets PersistentCheck field to given value.

### HasPersistentCheck

`func (o *Service) HasPersistentCheck() bool`

HasPersistentCheck returns a boolean if a field has been set.

### GetAgentProtocol

`func (o *Service) GetAgentProtocol() string`

GetAgentProtocol returns the AgentProtocol field if non-nil, zero value otherwise.

### GetAgentProtocolOk

`func (o *Service) GetAgentProtocolOk() (*string, bool)`

GetAgentProtocolOk returns a tuple with the AgentProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentProtocol

`func (o *Service) SetAgentProtocol(v string)`

SetAgentProtocol sets AgentProtocol field to given value.

### HasAgentProtocol

`func (o *Service) HasAgentProtocol() bool`

HasAgentProtocol returns a boolean if a field has been set.

### GetAgentInfo

`func (o *Service) GetAgentInfo() map[string]interface{}`

GetAgentInfo returns the AgentInfo field if non-nil, zero value otherwise.

### GetAgentInfoOk

`func (o *Service) GetAgentInfoOk() (*map[string]interface{}, bool)`

GetAgentInfoOk returns a tuple with the AgentInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentInfo

`func (o *Service) SetAgentInfo(v map[string]interface{})`

SetAgentInfo sets AgentInfo field to given value.

### HasAgentInfo

`func (o *Service) HasAgentInfo() bool`

HasAgentInfo returns a boolean if a field has been set.

### GetAgentInfoUrl

`func (o *Service) GetAgentInfoUrl() string`

GetAgentInfoUrl returns the AgentInfoUrl field if non-nil, zero value otherwise.

### GetAgentInfoUrlOk

`func (o *Service) GetAgentInfoUrlOk() (*string, bool)`

GetAgentInfoUrlOk returns a tuple with the AgentInfoUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentInfoUrl

`func (o *Service) SetAgentInfoUrl(v string)`

SetAgentInfoUrl sets AgentInfoUrl field to given value.

### HasAgentInfoUrl

`func (o *Service) HasAgentInfoUrl() bool`

HasAgentInfoUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


