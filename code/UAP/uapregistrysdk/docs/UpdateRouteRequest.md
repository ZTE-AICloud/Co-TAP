# UpdateRouteRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
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

### NewUpdateRouteRequest

`func NewUpdateRouteRequest(name string, ) *UpdateRouteRequest`

NewUpdateRouteRequest instantiates a new UpdateRouteRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateRouteRequestWithDefaults

`func NewUpdateRouteRequestWithDefaults() *UpdateRouteRequest`

NewUpdateRouteRequestWithDefaults instantiates a new UpdateRouteRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *UpdateRouteRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateRouteRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateRouteRequest) SetName(v string)`

SetName sets Name field to given value.


### GetGatewayId

`func (o *UpdateRouteRequest) GetGatewayId() string`

GetGatewayId returns the GatewayId field if non-nil, zero value otherwise.

### GetGatewayIdOk

`func (o *UpdateRouteRequest) GetGatewayIdOk() (*string, bool)`

GetGatewayIdOk returns a tuple with the GatewayId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGatewayId

`func (o *UpdateRouteRequest) SetGatewayId(v string)`

SetGatewayId sets GatewayId field to given value.

### HasGatewayId

`func (o *UpdateRouteRequest) HasGatewayId() bool`

HasGatewayId returns a boolean if a field has been set.

### GetProtocols

`func (o *UpdateRouteRequest) GetProtocols() []string`

GetProtocols returns the Protocols field if non-nil, zero value otherwise.

### GetProtocolsOk

`func (o *UpdateRouteRequest) GetProtocolsOk() (*[]string, bool)`

GetProtocolsOk returns a tuple with the Protocols field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocols

`func (o *UpdateRouteRequest) SetProtocols(v []string)`

SetProtocols sets Protocols field to given value.

### HasProtocols

`func (o *UpdateRouteRequest) HasProtocols() bool`

HasProtocols returns a boolean if a field has been set.

### GetMethods

`func (o *UpdateRouteRequest) GetMethods() []string`

GetMethods returns the Methods field if non-nil, zero value otherwise.

### GetMethodsOk

`func (o *UpdateRouteRequest) GetMethodsOk() (*[]string, bool)`

GetMethodsOk returns a tuple with the Methods field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMethods

`func (o *UpdateRouteRequest) SetMethods(v []string)`

SetMethods sets Methods field to given value.

### HasMethods

`func (o *UpdateRouteRequest) HasMethods() bool`

HasMethods returns a boolean if a field has been set.

### GetHosts

`func (o *UpdateRouteRequest) GetHosts() []string`

GetHosts returns the Hosts field if non-nil, zero value otherwise.

### GetHostsOk

`func (o *UpdateRouteRequest) GetHostsOk() (*[]string, bool)`

GetHostsOk returns a tuple with the Hosts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHosts

`func (o *UpdateRouteRequest) SetHosts(v []string)`

SetHosts sets Hosts field to given value.

### HasHosts

`func (o *UpdateRouteRequest) HasHosts() bool`

HasHosts returns a boolean if a field has been set.

### GetPaths

`func (o *UpdateRouteRequest) GetPaths() []string`

GetPaths returns the Paths field if non-nil, zero value otherwise.

### GetPathsOk

`func (o *UpdateRouteRequest) GetPathsOk() (*[]string, bool)`

GetPathsOk returns a tuple with the Paths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaths

`func (o *UpdateRouteRequest) SetPaths(v []string)`

SetPaths sets Paths field to given value.

### HasPaths

`func (o *UpdateRouteRequest) HasPaths() bool`

HasPaths returns a boolean if a field has been set.

### GetHeaders

`func (o *UpdateRouteRequest) GetHeaders() []string`

GetHeaders returns the Headers field if non-nil, zero value otherwise.

### GetHeadersOk

`func (o *UpdateRouteRequest) GetHeadersOk() (*[]string, bool)`

GetHeadersOk returns a tuple with the Headers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeaders

`func (o *UpdateRouteRequest) SetHeaders(v []string)`

SetHeaders sets Headers field to given value.

### HasHeaders

`func (o *UpdateRouteRequest) HasHeaders() bool`

HasHeaders returns a boolean if a field has been set.

### GetHttpsRedirectStatusCode

`func (o *UpdateRouteRequest) GetHttpsRedirectStatusCode() int32`

GetHttpsRedirectStatusCode returns the HttpsRedirectStatusCode field if non-nil, zero value otherwise.

### GetHttpsRedirectStatusCodeOk

`func (o *UpdateRouteRequest) GetHttpsRedirectStatusCodeOk() (*int32, bool)`

GetHttpsRedirectStatusCodeOk returns a tuple with the HttpsRedirectStatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHttpsRedirectStatusCode

`func (o *UpdateRouteRequest) SetHttpsRedirectStatusCode(v int32)`

SetHttpsRedirectStatusCode sets HttpsRedirectStatusCode field to given value.

### HasHttpsRedirectStatusCode

`func (o *UpdateRouteRequest) HasHttpsRedirectStatusCode() bool`

HasHttpsRedirectStatusCode returns a boolean if a field has been set.

### GetRegexPriority

`func (o *UpdateRouteRequest) GetRegexPriority() int32`

GetRegexPriority returns the RegexPriority field if non-nil, zero value otherwise.

### GetRegexPriorityOk

`func (o *UpdateRouteRequest) GetRegexPriorityOk() (*int32, bool)`

GetRegexPriorityOk returns a tuple with the RegexPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegexPriority

`func (o *UpdateRouteRequest) SetRegexPriority(v int32)`

SetRegexPriority sets RegexPriority field to given value.

### HasRegexPriority

`func (o *UpdateRouteRequest) HasRegexPriority() bool`

HasRegexPriority returns a boolean if a field has been set.

### GetStripPath

`func (o *UpdateRouteRequest) GetStripPath() bool`

GetStripPath returns the StripPath field if non-nil, zero value otherwise.

### GetStripPathOk

`func (o *UpdateRouteRequest) GetStripPathOk() (*bool, bool)`

GetStripPathOk returns a tuple with the StripPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripPath

`func (o *UpdateRouteRequest) SetStripPath(v bool)`

SetStripPath sets StripPath field to given value.

### HasStripPath

`func (o *UpdateRouteRequest) HasStripPath() bool`

HasStripPath returns a boolean if a field has been set.

### GetPreserveHost

`func (o *UpdateRouteRequest) GetPreserveHost() bool`

GetPreserveHost returns the PreserveHost field if non-nil, zero value otherwise.

### GetPreserveHostOk

`func (o *UpdateRouteRequest) GetPreserveHostOk() (*bool, bool)`

GetPreserveHostOk returns a tuple with the PreserveHost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveHost

`func (o *UpdateRouteRequest) SetPreserveHost(v bool)`

SetPreserveHost sets PreserveHost field to given value.

### HasPreserveHost

`func (o *UpdateRouteRequest) HasPreserveHost() bool`

HasPreserveHost returns a boolean if a field has been set.

### GetRequestBuffering

`func (o *UpdateRouteRequest) GetRequestBuffering() bool`

GetRequestBuffering returns the RequestBuffering field if non-nil, zero value otherwise.

### GetRequestBufferingOk

`func (o *UpdateRouteRequest) GetRequestBufferingOk() (*bool, bool)`

GetRequestBufferingOk returns a tuple with the RequestBuffering field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestBuffering

`func (o *UpdateRouteRequest) SetRequestBuffering(v bool)`

SetRequestBuffering sets RequestBuffering field to given value.

### HasRequestBuffering

`func (o *UpdateRouteRequest) HasRequestBuffering() bool`

HasRequestBuffering returns a boolean if a field has been set.

### GetResponseBuffering

`func (o *UpdateRouteRequest) GetResponseBuffering() bool`

GetResponseBuffering returns the ResponseBuffering field if non-nil, zero value otherwise.

### GetResponseBufferingOk

`func (o *UpdateRouteRequest) GetResponseBufferingOk() (*bool, bool)`

GetResponseBufferingOk returns a tuple with the ResponseBuffering field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseBuffering

`func (o *UpdateRouteRequest) SetResponseBuffering(v bool)`

SetResponseBuffering sets ResponseBuffering field to given value.

### HasResponseBuffering

`func (o *UpdateRouteRequest) HasResponseBuffering() bool`

HasResponseBuffering returns a boolean if a field has been set.

### GetSnis

`func (o *UpdateRouteRequest) GetSnis() []string`

GetSnis returns the Snis field if non-nil, zero value otherwise.

### GetSnisOk

`func (o *UpdateRouteRequest) GetSnisOk() (*[]string, bool)`

GetSnisOk returns a tuple with the Snis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnis

`func (o *UpdateRouteRequest) SetSnis(v []string)`

SetSnis sets Snis field to given value.

### HasSnis

`func (o *UpdateRouteRequest) HasSnis() bool`

HasSnis returns a boolean if a field has been set.

### GetSources

`func (o *UpdateRouteRequest) GetSources() []string`

GetSources returns the Sources field if non-nil, zero value otherwise.

### GetSourcesOk

`func (o *UpdateRouteRequest) GetSourcesOk() (*[]string, bool)`

GetSourcesOk returns a tuple with the Sources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSources

`func (o *UpdateRouteRequest) SetSources(v []string)`

SetSources sets Sources field to given value.

### HasSources

`func (o *UpdateRouteRequest) HasSources() bool`

HasSources returns a boolean if a field has been set.

### GetDestinations

`func (o *UpdateRouteRequest) GetDestinations() []string`

GetDestinations returns the Destinations field if non-nil, zero value otherwise.

### GetDestinationsOk

`func (o *UpdateRouteRequest) GetDestinationsOk() (*[]string, bool)`

GetDestinationsOk returns a tuple with the Destinations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinations

`func (o *UpdateRouteRequest) SetDestinations(v []string)`

SetDestinations sets Destinations field to given value.

### HasDestinations

`func (o *UpdateRouteRequest) HasDestinations() bool`

HasDestinations returns a boolean if a field has been set.

### GetTags

`func (o *UpdateRouteRequest) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *UpdateRouteRequest) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *UpdateRouteRequest) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *UpdateRouteRequest) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetAgentProtocol

`func (o *UpdateRouteRequest) GetAgentProtocol() string`

GetAgentProtocol returns the AgentProtocol field if non-nil, zero value otherwise.

### GetAgentProtocolOk

`func (o *UpdateRouteRequest) GetAgentProtocolOk() (*string, bool)`

GetAgentProtocolOk returns a tuple with the AgentProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentProtocol

`func (o *UpdateRouteRequest) SetAgentProtocol(v string)`

SetAgentProtocol sets AgentProtocol field to given value.

### HasAgentProtocol

`func (o *UpdateRouteRequest) HasAgentProtocol() bool`

HasAgentProtocol returns a boolean if a field has been set.

### GetService

`func (o *UpdateRouteRequest) GetService() string`

GetService returns the Service field if non-nil, zero value otherwise.

### GetServiceOk

`func (o *UpdateRouteRequest) GetServiceOk() (*string, bool)`

GetServiceOk returns a tuple with the Service field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetService

`func (o *UpdateRouteRequest) SetService(v string)`

SetService sets Service field to given value.

### HasService

`func (o *UpdateRouteRequest) HasService() bool`

HasService returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


