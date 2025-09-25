# Route

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | 路由唯一标识 | [optional] 
**CreatedAt** | Pointer to **int64** | 记录创建时间（由注册中心设置） | [optional] 
**UpdatedAt** | Pointer to **int64** | 记录最后更新时间（由注册中心设置） | [optional] 
**Index** | Pointer to **int64** | 路由更新序号 | [optional] 
**Name** | **string** | 路由名称 | 
**GatewayId** | Pointer to **string** | 此路由关联的网关标识 | [optional] 
**Protocols** | Pointer to **[]string** | 支持的协议列表 | [optional] 
**Methods** | Pointer to **[]string** | 匹配的HTTP方法列表 | [optional] 
**Hosts** | Pointer to **[]string** | 匹配的Host列表 | [optional] 
**Paths** | Pointer to **[]string** | 匹配的路径 | [optional] 
**Headers** | Pointer to **[]string** | 匹配的请求头 | [optional] 
**HttpsRedirectStatusCode** | Pointer to **int32** | 协议不匹配时的重定向状态码 | [optional] 
**RegexPriority** | Pointer to **int32** | 正则匹配优先级 | [optional] 
**StripPath** | Pointer to **bool** | 是否移除路径前缀 | [optional] 
**PreserveHost** | Pointer to **bool** | 是否保留原始请求的Host头 | [optional] 
**RequestBuffering** | Pointer to **bool** | 是否缓冲请求体数据 | [optional] 
**ResponseBuffering** | Pointer to **bool** | 是否缓冲响应体数据 | [optional] 
**Snis** | Pointer to **[]string** | 支持的SNI列表 | [optional] 
**Sources** | Pointer to **[]string** | 允许的源IP列表 | [optional] 
**Destinations** | Pointer to **[]string** | 目标服务地址列表 | [optional] 
**Tags** | Pointer to **[]string** | 路由标签 | [optional] 
**AgentProtocol** | Pointer to **string** | Agent/Tool通信协议 | [optional] 
**Service** | Pointer to **string** | 关联的服务名 | [optional] 

## Methods

### NewRoute

`func NewRoute(name string, ) *Route`

NewRoute instantiates a new Route object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRouteWithDefaults

`func NewRouteWithDefaults() *Route`

NewRouteWithDefaults instantiates a new Route object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Route) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Route) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Route) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Route) HasId() bool`

HasId returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Route) GetCreatedAt() int64`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Route) GetCreatedAtOk() (*int64, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Route) SetCreatedAt(v int64)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Route) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *Route) GetUpdatedAt() int64`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Route) GetUpdatedAtOk() (*int64, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Route) SetUpdatedAt(v int64)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *Route) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetIndex

`func (o *Route) GetIndex() int64`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *Route) GetIndexOk() (*int64, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *Route) SetIndex(v int64)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *Route) HasIndex() bool`

HasIndex returns a boolean if a field has been set.

### GetName

`func (o *Route) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Route) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Route) SetName(v string)`

SetName sets Name field to given value.


### GetGatewayId

`func (o *Route) GetGatewayId() string`

GetGatewayId returns the GatewayId field if non-nil, zero value otherwise.

### GetGatewayIdOk

`func (o *Route) GetGatewayIdOk() (*string, bool)`

GetGatewayIdOk returns a tuple with the GatewayId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGatewayId

`func (o *Route) SetGatewayId(v string)`

SetGatewayId sets GatewayId field to given value.

### HasGatewayId

`func (o *Route) HasGatewayId() bool`

HasGatewayId returns a boolean if a field has been set.

### GetProtocols

`func (o *Route) GetProtocols() []string`

GetProtocols returns the Protocols field if non-nil, zero value otherwise.

### GetProtocolsOk

`func (o *Route) GetProtocolsOk() (*[]string, bool)`

GetProtocolsOk returns a tuple with the Protocols field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocols

`func (o *Route) SetProtocols(v []string)`

SetProtocols sets Protocols field to given value.

### HasProtocols

`func (o *Route) HasProtocols() bool`

HasProtocols returns a boolean if a field has been set.

### GetMethods

`func (o *Route) GetMethods() []string`

GetMethods returns the Methods field if non-nil, zero value otherwise.

### GetMethodsOk

`func (o *Route) GetMethodsOk() (*[]string, bool)`

GetMethodsOk returns a tuple with the Methods field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMethods

`func (o *Route) SetMethods(v []string)`

SetMethods sets Methods field to given value.

### HasMethods

`func (o *Route) HasMethods() bool`

HasMethods returns a boolean if a field has been set.

### GetHosts

`func (o *Route) GetHosts() []string`

GetHosts returns the Hosts field if non-nil, zero value otherwise.

### GetHostsOk

`func (o *Route) GetHostsOk() (*[]string, bool)`

GetHostsOk returns a tuple with the Hosts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHosts

`func (o *Route) SetHosts(v []string)`

SetHosts sets Hosts field to given value.

### HasHosts

`func (o *Route) HasHosts() bool`

HasHosts returns a boolean if a field has been set.

### GetPaths

`func (o *Route) GetPaths() []string`

GetPaths returns the Paths field if non-nil, zero value otherwise.

### GetPathsOk

`func (o *Route) GetPathsOk() (*[]string, bool)`

GetPathsOk returns a tuple with the Paths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaths

`func (o *Route) SetPaths(v []string)`

SetPaths sets Paths field to given value.

### HasPaths

`func (o *Route) HasPaths() bool`

HasPaths returns a boolean if a field has been set.

### GetHeaders

`func (o *Route) GetHeaders() []string`

GetHeaders returns the Headers field if non-nil, zero value otherwise.

### GetHeadersOk

`func (o *Route) GetHeadersOk() (*[]string, bool)`

GetHeadersOk returns a tuple with the Headers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeaders

`func (o *Route) SetHeaders(v []string)`

SetHeaders sets Headers field to given value.

### HasHeaders

`func (o *Route) HasHeaders() bool`

HasHeaders returns a boolean if a field has been set.

### GetHttpsRedirectStatusCode

`func (o *Route) GetHttpsRedirectStatusCode() int32`

GetHttpsRedirectStatusCode returns the HttpsRedirectStatusCode field if non-nil, zero value otherwise.

### GetHttpsRedirectStatusCodeOk

`func (o *Route) GetHttpsRedirectStatusCodeOk() (*int32, bool)`

GetHttpsRedirectStatusCodeOk returns a tuple with the HttpsRedirectStatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHttpsRedirectStatusCode

`func (o *Route) SetHttpsRedirectStatusCode(v int32)`

SetHttpsRedirectStatusCode sets HttpsRedirectStatusCode field to given value.

### HasHttpsRedirectStatusCode

`func (o *Route) HasHttpsRedirectStatusCode() bool`

HasHttpsRedirectStatusCode returns a boolean if a field has been set.

### GetRegexPriority

`func (o *Route) GetRegexPriority() int32`

GetRegexPriority returns the RegexPriority field if non-nil, zero value otherwise.

### GetRegexPriorityOk

`func (o *Route) GetRegexPriorityOk() (*int32, bool)`

GetRegexPriorityOk returns a tuple with the RegexPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegexPriority

`func (o *Route) SetRegexPriority(v int32)`

SetRegexPriority sets RegexPriority field to given value.

### HasRegexPriority

`func (o *Route) HasRegexPriority() bool`

HasRegexPriority returns a boolean if a field has been set.

### GetStripPath

`func (o *Route) GetStripPath() bool`

GetStripPath returns the StripPath field if non-nil, zero value otherwise.

### GetStripPathOk

`func (o *Route) GetStripPathOk() (*bool, bool)`

GetStripPathOk returns a tuple with the StripPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripPath

`func (o *Route) SetStripPath(v bool)`

SetStripPath sets StripPath field to given value.

### HasStripPath

`func (o *Route) HasStripPath() bool`

HasStripPath returns a boolean if a field has been set.

### GetPreserveHost

`func (o *Route) GetPreserveHost() bool`

GetPreserveHost returns the PreserveHost field if non-nil, zero value otherwise.

### GetPreserveHostOk

`func (o *Route) GetPreserveHostOk() (*bool, bool)`

GetPreserveHostOk returns a tuple with the PreserveHost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveHost

`func (o *Route) SetPreserveHost(v bool)`

SetPreserveHost sets PreserveHost field to given value.

### HasPreserveHost

`func (o *Route) HasPreserveHost() bool`

HasPreserveHost returns a boolean if a field has been set.

### GetRequestBuffering

`func (o *Route) GetRequestBuffering() bool`

GetRequestBuffering returns the RequestBuffering field if non-nil, zero value otherwise.

### GetRequestBufferingOk

`func (o *Route) GetRequestBufferingOk() (*bool, bool)`

GetRequestBufferingOk returns a tuple with the RequestBuffering field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestBuffering

`func (o *Route) SetRequestBuffering(v bool)`

SetRequestBuffering sets RequestBuffering field to given value.

### HasRequestBuffering

`func (o *Route) HasRequestBuffering() bool`

HasRequestBuffering returns a boolean if a field has been set.

### GetResponseBuffering

`func (o *Route) GetResponseBuffering() bool`

GetResponseBuffering returns the ResponseBuffering field if non-nil, zero value otherwise.

### GetResponseBufferingOk

`func (o *Route) GetResponseBufferingOk() (*bool, bool)`

GetResponseBufferingOk returns a tuple with the ResponseBuffering field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseBuffering

`func (o *Route) SetResponseBuffering(v bool)`

SetResponseBuffering sets ResponseBuffering field to given value.

### HasResponseBuffering

`func (o *Route) HasResponseBuffering() bool`

HasResponseBuffering returns a boolean if a field has been set.

### GetSnis

`func (o *Route) GetSnis() []string`

GetSnis returns the Snis field if non-nil, zero value otherwise.

### GetSnisOk

`func (o *Route) GetSnisOk() (*[]string, bool)`

GetSnisOk returns a tuple with the Snis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnis

`func (o *Route) SetSnis(v []string)`

SetSnis sets Snis field to given value.

### HasSnis

`func (o *Route) HasSnis() bool`

HasSnis returns a boolean if a field has been set.

### GetSources

`func (o *Route) GetSources() []string`

GetSources returns the Sources field if non-nil, zero value otherwise.

### GetSourcesOk

`func (o *Route) GetSourcesOk() (*[]string, bool)`

GetSourcesOk returns a tuple with the Sources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSources

`func (o *Route) SetSources(v []string)`

SetSources sets Sources field to given value.

### HasSources

`func (o *Route) HasSources() bool`

HasSources returns a boolean if a field has been set.

### GetDestinations

`func (o *Route) GetDestinations() []string`

GetDestinations returns the Destinations field if non-nil, zero value otherwise.

### GetDestinationsOk

`func (o *Route) GetDestinationsOk() (*[]string, bool)`

GetDestinationsOk returns a tuple with the Destinations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinations

`func (o *Route) SetDestinations(v []string)`

SetDestinations sets Destinations field to given value.

### HasDestinations

`func (o *Route) HasDestinations() bool`

HasDestinations returns a boolean if a field has been set.

### GetTags

`func (o *Route) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *Route) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *Route) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *Route) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetAgentProtocol

`func (o *Route) GetAgentProtocol() string`

GetAgentProtocol returns the AgentProtocol field if non-nil, zero value otherwise.

### GetAgentProtocolOk

`func (o *Route) GetAgentProtocolOk() (*string, bool)`

GetAgentProtocolOk returns a tuple with the AgentProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentProtocol

`func (o *Route) SetAgentProtocol(v string)`

SetAgentProtocol sets AgentProtocol field to given value.

### HasAgentProtocol

`func (o *Route) HasAgentProtocol() bool`

HasAgentProtocol returns a boolean if a field has been set.

### GetService

`func (o *Route) GetService() string`

GetService returns the Service field if non-nil, zero value otherwise.

### GetServiceOk

`func (o *Route) GetServiceOk() (*string, bool)`

GetServiceOk returns a tuple with the Service field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetService

`func (o *Route) SetService(v string)`

SetService sets Service field to given value.

### HasService

`func (o *Route) HasService() bool`

HasService returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


