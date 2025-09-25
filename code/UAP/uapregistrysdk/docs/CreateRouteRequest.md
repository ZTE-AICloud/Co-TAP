# CreateRouteRequest

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

### NewCreateRouteRequest

`func NewCreateRouteRequest(name string, ) *CreateRouteRequest`

NewCreateRouteRequest instantiates a new CreateRouteRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateRouteRequestWithDefaults

`func NewCreateRouteRequestWithDefaults() *CreateRouteRequest`

NewCreateRouteRequestWithDefaults instantiates a new CreateRouteRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CreateRouteRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateRouteRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateRouteRequest) SetName(v string)`

SetName sets Name field to given value.


### GetGatewayId

`func (o *CreateRouteRequest) GetGatewayId() string`

GetGatewayId returns the GatewayId field if non-nil, zero value otherwise.

### GetGatewayIdOk

`func (o *CreateRouteRequest) GetGatewayIdOk() (*string, bool)`

GetGatewayIdOk returns a tuple with the GatewayId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGatewayId

`func (o *CreateRouteRequest) SetGatewayId(v string)`

SetGatewayId sets GatewayId field to given value.

### HasGatewayId

`func (o *CreateRouteRequest) HasGatewayId() bool`

HasGatewayId returns a boolean if a field has been set.

### GetProtocols

`func (o *CreateRouteRequest) GetProtocols() []string`

GetProtocols returns the Protocols field if non-nil, zero value otherwise.

### GetProtocolsOk

`func (o *CreateRouteRequest) GetProtocolsOk() (*[]string, bool)`

GetProtocolsOk returns a tuple with the Protocols field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocols

`func (o *CreateRouteRequest) SetProtocols(v []string)`

SetProtocols sets Protocols field to given value.

### HasProtocols

`func (o *CreateRouteRequest) HasProtocols() bool`

HasProtocols returns a boolean if a field has been set.

### GetMethods

`func (o *CreateRouteRequest) GetMethods() []string`

GetMethods returns the Methods field if non-nil, zero value otherwise.

### GetMethodsOk

`func (o *CreateRouteRequest) GetMethodsOk() (*[]string, bool)`

GetMethodsOk returns a tuple with the Methods field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMethods

`func (o *CreateRouteRequest) SetMethods(v []string)`

SetMethods sets Methods field to given value.

### HasMethods

`func (o *CreateRouteRequest) HasMethods() bool`

HasMethods returns a boolean if a field has been set.

### GetHosts

`func (o *CreateRouteRequest) GetHosts() []string`

GetHosts returns the Hosts field if non-nil, zero value otherwise.

### GetHostsOk

`func (o *CreateRouteRequest) GetHostsOk() (*[]string, bool)`

GetHostsOk returns a tuple with the Hosts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHosts

`func (o *CreateRouteRequest) SetHosts(v []string)`

SetHosts sets Hosts field to given value.

### HasHosts

`func (o *CreateRouteRequest) HasHosts() bool`

HasHosts returns a boolean if a field has been set.

### GetPaths

`func (o *CreateRouteRequest) GetPaths() []string`

GetPaths returns the Paths field if non-nil, zero value otherwise.

### GetPathsOk

`func (o *CreateRouteRequest) GetPathsOk() (*[]string, bool)`

GetPathsOk returns a tuple with the Paths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaths

`func (o *CreateRouteRequest) SetPaths(v []string)`

SetPaths sets Paths field to given value.

### HasPaths

`func (o *CreateRouteRequest) HasPaths() bool`

HasPaths returns a boolean if a field has been set.

### GetHeaders

`func (o *CreateRouteRequest) GetHeaders() []string`

GetHeaders returns the Headers field if non-nil, zero value otherwise.

### GetHeadersOk

`func (o *CreateRouteRequest) GetHeadersOk() (*[]string, bool)`

GetHeadersOk returns a tuple with the Headers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeaders

`func (o *CreateRouteRequest) SetHeaders(v []string)`

SetHeaders sets Headers field to given value.

### HasHeaders

`func (o *CreateRouteRequest) HasHeaders() bool`

HasHeaders returns a boolean if a field has been set.

### GetHttpsRedirectStatusCode

`func (o *CreateRouteRequest) GetHttpsRedirectStatusCode() int32`

GetHttpsRedirectStatusCode returns the HttpsRedirectStatusCode field if non-nil, zero value otherwise.

### GetHttpsRedirectStatusCodeOk

`func (o *CreateRouteRequest) GetHttpsRedirectStatusCodeOk() (*int32, bool)`

GetHttpsRedirectStatusCodeOk returns a tuple with the HttpsRedirectStatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHttpsRedirectStatusCode

`func (o *CreateRouteRequest) SetHttpsRedirectStatusCode(v int32)`

SetHttpsRedirectStatusCode sets HttpsRedirectStatusCode field to given value.

### HasHttpsRedirectStatusCode

`func (o *CreateRouteRequest) HasHttpsRedirectStatusCode() bool`

HasHttpsRedirectStatusCode returns a boolean if a field has been set.

### GetRegexPriority

`func (o *CreateRouteRequest) GetRegexPriority() int32`

GetRegexPriority returns the RegexPriority field if non-nil, zero value otherwise.

### GetRegexPriorityOk

`func (o *CreateRouteRequest) GetRegexPriorityOk() (*int32, bool)`

GetRegexPriorityOk returns a tuple with the RegexPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegexPriority

`func (o *CreateRouteRequest) SetRegexPriority(v int32)`

SetRegexPriority sets RegexPriority field to given value.

### HasRegexPriority

`func (o *CreateRouteRequest) HasRegexPriority() bool`

HasRegexPriority returns a boolean if a field has been set.

### GetStripPath

`func (o *CreateRouteRequest) GetStripPath() bool`

GetStripPath returns the StripPath field if non-nil, zero value otherwise.

### GetStripPathOk

`func (o *CreateRouteRequest) GetStripPathOk() (*bool, bool)`

GetStripPathOk returns a tuple with the StripPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStripPath

`func (o *CreateRouteRequest) SetStripPath(v bool)`

SetStripPath sets StripPath field to given value.

### HasStripPath

`func (o *CreateRouteRequest) HasStripPath() bool`

HasStripPath returns a boolean if a field has been set.

### GetPreserveHost

`func (o *CreateRouteRequest) GetPreserveHost() bool`

GetPreserveHost returns the PreserveHost field if non-nil, zero value otherwise.

### GetPreserveHostOk

`func (o *CreateRouteRequest) GetPreserveHostOk() (*bool, bool)`

GetPreserveHostOk returns a tuple with the PreserveHost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreserveHost

`func (o *CreateRouteRequest) SetPreserveHost(v bool)`

SetPreserveHost sets PreserveHost field to given value.

### HasPreserveHost

`func (o *CreateRouteRequest) HasPreserveHost() bool`

HasPreserveHost returns a boolean if a field has been set.

### GetRequestBuffering

`func (o *CreateRouteRequest) GetRequestBuffering() bool`

GetRequestBuffering returns the RequestBuffering field if non-nil, zero value otherwise.

### GetRequestBufferingOk

`func (o *CreateRouteRequest) GetRequestBufferingOk() (*bool, bool)`

GetRequestBufferingOk returns a tuple with the RequestBuffering field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestBuffering

`func (o *CreateRouteRequest) SetRequestBuffering(v bool)`

SetRequestBuffering sets RequestBuffering field to given value.

### HasRequestBuffering

`func (o *CreateRouteRequest) HasRequestBuffering() bool`

HasRequestBuffering returns a boolean if a field has been set.

### GetResponseBuffering

`func (o *CreateRouteRequest) GetResponseBuffering() bool`

GetResponseBuffering returns the ResponseBuffering field if non-nil, zero value otherwise.

### GetResponseBufferingOk

`func (o *CreateRouteRequest) GetResponseBufferingOk() (*bool, bool)`

GetResponseBufferingOk returns a tuple with the ResponseBuffering field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseBuffering

`func (o *CreateRouteRequest) SetResponseBuffering(v bool)`

SetResponseBuffering sets ResponseBuffering field to given value.

### HasResponseBuffering

`func (o *CreateRouteRequest) HasResponseBuffering() bool`

HasResponseBuffering returns a boolean if a field has been set.

### GetSnis

`func (o *CreateRouteRequest) GetSnis() []string`

GetSnis returns the Snis field if non-nil, zero value otherwise.

### GetSnisOk

`func (o *CreateRouteRequest) GetSnisOk() (*[]string, bool)`

GetSnisOk returns a tuple with the Snis field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnis

`func (o *CreateRouteRequest) SetSnis(v []string)`

SetSnis sets Snis field to given value.

### HasSnis

`func (o *CreateRouteRequest) HasSnis() bool`

HasSnis returns a boolean if a field has been set.

### GetSources

`func (o *CreateRouteRequest) GetSources() []string`

GetSources returns the Sources field if non-nil, zero value otherwise.

### GetSourcesOk

`func (o *CreateRouteRequest) GetSourcesOk() (*[]string, bool)`

GetSourcesOk returns a tuple with the Sources field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSources

`func (o *CreateRouteRequest) SetSources(v []string)`

SetSources sets Sources field to given value.

### HasSources

`func (o *CreateRouteRequest) HasSources() bool`

HasSources returns a boolean if a field has been set.

### GetDestinations

`func (o *CreateRouteRequest) GetDestinations() []string`

GetDestinations returns the Destinations field if non-nil, zero value otherwise.

### GetDestinationsOk

`func (o *CreateRouteRequest) GetDestinationsOk() (*[]string, bool)`

GetDestinationsOk returns a tuple with the Destinations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinations

`func (o *CreateRouteRequest) SetDestinations(v []string)`

SetDestinations sets Destinations field to given value.

### HasDestinations

`func (o *CreateRouteRequest) HasDestinations() bool`

HasDestinations returns a boolean if a field has been set.

### GetTags

`func (o *CreateRouteRequest) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *CreateRouteRequest) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *CreateRouteRequest) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *CreateRouteRequest) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetAgentProtocol

`func (o *CreateRouteRequest) GetAgentProtocol() string`

GetAgentProtocol returns the AgentProtocol field if non-nil, zero value otherwise.

### GetAgentProtocolOk

`func (o *CreateRouteRequest) GetAgentProtocolOk() (*string, bool)`

GetAgentProtocolOk returns a tuple with the AgentProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgentProtocol

`func (o *CreateRouteRequest) SetAgentProtocol(v string)`

SetAgentProtocol sets AgentProtocol field to given value.

### HasAgentProtocol

`func (o *CreateRouteRequest) HasAgentProtocol() bool`

HasAgentProtocol returns a boolean if a field has been set.

### GetService

`func (o *CreateRouteRequest) GetService() string`

GetService returns the Service field if non-nil, zero value otherwise.

### GetServiceOk

`func (o *CreateRouteRequest) GetServiceOk() (*string, bool)`

GetServiceOk returns a tuple with the Service field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetService

`func (o *CreateRouteRequest) SetService(v string)`

SetService sets Service field to given value.

### HasService

`func (o *CreateRouteRequest) HasService() bool`

HasService returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


