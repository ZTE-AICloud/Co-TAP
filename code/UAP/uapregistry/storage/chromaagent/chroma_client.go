package agent

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/amikos-tech/chroma-go/pkg/embeddings"
)

const DefaultSimilarityThreshold float32 = 0.5

// 预编译正则，避免每次调用重复编译，提升性能
var (
	// 匹配需要移除的无意义特殊符号（保留中英文常用标点、汉字、数字、字母、空格）
	specialCharRegex = regexp.MustCompile(`[^a-zA-Z0-9\u4e00-\u9fa5，。？！；：、“”‘’（）【】《》……—,.?!;:'"()\[\] ]`)
	// 匹配连续重复的标点符号，合并为单个
	repeatPunctRegex = regexp.MustCompile(`([，。？！；：、,.?!;:])\1+`)
)

// NewChromaClient 初始化 Agent
func NewChromaClient(timeout time.Duration) (*ChromaClient, error) {
	httpClient := &http.Client{
		Timeout: timeout,
	}

	client, err := chroma.NewHTTPClient(chroma.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &ChromaClient{
		client:     client,
		collection: make(map[string]chroma.Collection),
	}, nil
}

// CreateCollection 创建或获取集合并缓存
func (c *ChromaClient) CreateCollection(ctx context.Context, name string, embeddingFunc embeddings.EmbeddingFunction, metadata map[string]interface{}) error {
	var opts []chroma.CreateCollectionOption

	if embeddingFunc != nil {
		opts = append(opts, chroma.WithEmbeddingFunctionCreate(embeddingFunc))
	}

	if metadata != nil {
		colMeta := chroma.NewMetadataFromMap(metadata)
		opts = append(opts, chroma.WithCollectionMetadataCreate(colMeta))
	}

	col, err := c.client.GetOrCreateCollection(ctx, name, opts...)
	if err != nil {
		return fmt.Errorf("create collection failed: %w", err)
	}

	c.Lock()
	c.collection[name] = col
	c.Unlock()

	return nil
}

// GetCollection 获取集合（优先读本地缓存，未命中则请求服务端）
func (c *ChromaClient) GetCollection(ctx context.Context, name string) (chroma.Collection, error) {
	// 1. 尝试从本地缓存读取（读锁）
	c.RLock()
	col, ok := c.collection[name]
	c.RUnlock()
	if ok {
		return col, nil
	}

	// 2. 缓存未命中，向 ChromaDB 服务端请求
	// 传入 nil 作为 embeddingFunc，Chroma 会自动使用服务端默认的嵌入函数
	col, err := c.client.GetCollection(ctx, name, nil)
	if err != nil {
		return nil, fmt.Errorf("get collection %s failed: %w", name, err)
	}

	// 3. 写入本地缓存（写锁）
	c.Lock()
	c.collection[name] = col
	c.Unlock()

	return col, nil
}

// AddDocumentsAutoID 批量添加文档到指定集合（自动生成ID）
func (c *ChromaClient) AddDocumentsAutoID(ctx context.Context, collectionName string, documents []string, metadatas []map[string]interface{}) error {
	// 边界校验：文档不能为空
	if len(documents) == 0 {
		return nil
	}

	col, err := c.GetCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	// 构建 Add 选项
	addOpts := []chroma.AddOption{
		chroma.WithTexts(documents...),
	}

	// 仅当 metadatas 非空且长度匹配时，才注入 metadata
	if len(metadatas) > 0 {
		if len(metadatas) != len(documents) {
			return fmt.Errorf("metadatas length %d does not match documents length %d", len(metadatas), len(documents))
		}

		docMetadatas := make([]chroma.DocumentMetadata, len(metadatas))
		for i, meta := range metadatas {
			// 根包工具函数：map 转 DocumentMetadata 接口实现
			docMeta := chroma.NewMetadataFromMap(meta)
			docMetadatas[i] = docMeta
		}
		addOpts = append(addOpts, chroma.WithMetadatas(docMetadatas...))
	}

	// 不传入 WithIDs，Chroma 自动生成 UUID
	err = col.Add(ctx, addOpts...)
	if err != nil {
		return fmt.Errorf("add documents failed: %w", err)
	}
	return nil
}

// UpsertSingleDocument 单条幂等写入（便捷方法）
func (c *ChromaClient) UpsertSingleDocument(
    ctx context.Context,
    collectionName string,
    id string,
    document string,
    metadata map[string]interface{},
) error {
    var metadatas []map[string]interface{}
    if metadata != nil {
        metadatas = []map[string]interface{}{metadata}
    }
    return c.UpsertDocuments(ctx, &ChromaUpsertRequest{
        CollectionName: collectionName,
        IDs:            []string{id},
        Documents:      []string{document},
        Metadatas:      metadatas,
    })
}

func (c *ChromaClient) UpsertDocuments(ctx context.Context, req *ChromaUpsertRequest) error {
	// 入参防御
	if req == nil {
		return fmt.Errorf("upsert request is nil")
	}

	// 基础参数校验
	if len(req.IDs) == 0 || len(req.Documents) == 0 {
		return nil
	}
	if len(req.IDs) != len(req.Documents) {
		return fmt.Errorf("ids length %d does not match documents length %d", len(req.IDs), len(req.Documents))
	}

	col, err := c.GetCollection(ctx, req.CollectionName)
	if err != nil {
		return err
	}

	docIDs := make([]chroma.DocumentID, len(req.IDs))
	for i, id := range req.IDs {
		docIDs[i] = chroma.DocumentID(id)
	}

	upsertOpts, err := buildUpsertOption(docIDs, req)
	if err != nil {
		return err
	}

	err = col.Upsert(ctx, upsertOpts...)
	if err != nil {
		return fmt.Errorf("upsert documents failed: %w", err)
	}
	return nil
}

// buildUpsertOption 构建 Upsert 参数，严格校验 metadata 格式
func buildUpsertOption(docIDs []chroma.DocumentID, req *ChromaUpsertRequest) ([]chroma.CollectionAddOption, error) {
	opts := []chroma.CollectionAddOption{
		chroma.WithIDs(docIDs...),
		chroma.WithTexts(req.Documents...),
	}

	if len(req.Metadatas) == 0 {
		return opts, nil
	}
	if len(req.Metadatas) != len(req.Documents) {
		return nil, fmt.Errorf("metadatas length %d does not match documents length %d", len(req.Metadatas), len(req.Documents))
	}

	// 严格模式转换 metadata：类型不合法直接报错，避免静默丢失字段
	docMetadatas := make([]chroma.DocumentMetadata, len(req.Metadatas))
	for i, meta := range req.Metadatas {
		// 适配数组类型：将 []string 转为 []interface{}，保证 Chroma 识别为数组字段
		normalizedMeta := normalizeMetadataArray(meta)

		docMeta, err := chroma.NewMetadataFromMapStrict(normalizedMeta)
		if err != nil {
			return nil, fmt.Errorf("metadata format error at index %d: %w", i, err)
		}
		docMetadatas[i] = docMeta
	}
	opts = append(opts, chroma.WithMetadatas(docMetadatas...))
	return opts, nil
}

// normalizeMetadataArray 标准化 metadata 中的数组类型
func normalizeMetadataArray(meta map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(meta))
	for k, v := range meta {
		switch val := v.(type) {
		case []string:
			// []string 转 []interface{}，保证 Chroma 识别为数组类型
			arr := make([]interface{}, len(val))
			for i, item := range val {
				arr[i] = item
			}
			result[k] = arr
		default:
			result[k] = v
		}
	}
	return result
}

// Query 语义查询：硬约束前置过滤 + 相似度阈值后处理
func (c *ChromaClient) Query(ctx context.Context, cqr *ChromaQueryRequest) (chroma.QueryResult, error) {
	col, err := c.GetCollection(ctx, cqr.CollectionName)
	if err != nil {
		return nil, err
	}

	// 1. 动态构建前置 where 过滤条件
	queryOpts := []chroma.QueryOption{
		chroma.WithQueryTexts(PreprocessQuery(cqr.Query)),
		chroma.WithNResults(cqr.Top_K),
	}

	whereFilter := buildWhereFilter(cqr)
	if whereFilter != nil {
		queryOpts = append(queryOpts, chroma.WithWhere(whereFilter))
	}

	// 2. 执行向量查询（Chroma 内部已完成硬约束过滤）
	res, err := col.Query(ctx, queryOpts...)
	if err != nil {
		return nil, fmt.Errorf("chroma query failed: %w", err)
	}

	filteredResult, err := queryFilter(res, cqr.Threshold)
	if err != nil {
		return nil, err
	}

	return filteredResult, nil
}

// DeleteDocuments 根据 ID 列表删除文档
func (c *ChromaClient) DeleteDocuments(ctx context.Context, collectionName string, ids []string) error {
	col, err := c.GetCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	docIDs := make([]chroma.DocumentID, len(ids))
	for i, id := range ids {
		docIDs[i] = chroma.DocumentID(id)
	}

	err = col.Delete(ctx,
		chroma.WithIDs(docIDs...), // 传入 DocumentID 类型的可变参数
	)
	if err != nil {
		return fmt.Errorf("delete documents failed: %w", err)
	}
	return nil
}

// DeleteCollection 删除集合并清理本地缓存
func (c *ChromaClient) DeleteCollection(ctx context.Context, name string) error {
	err := c.client.DeleteCollection(ctx, name)
	if err != nil {
		return fmt.Errorf("delete collection failed: %w", err)
	}

	// 清理本地缓存
	c.Lock()
	delete(c.collection, name)
	c.Unlock()

	return nil
}

// PreprocessQueryEnhanced 文本规范化：增强版
// 功能：空白压缩 + 英文小写 + 清理无效特殊符号 + 合并重复标点 + 保留中文语义标点
func PreprocessQuery(queryText string) string {
	// 1. 全角空格统一转为半角空格
	cleaned := strings.ReplaceAll(queryText, "　", " ")

	// 2. 压缩所有连续空白为单个空格，去除首尾空白
	fields := strings.Fields(cleaned)
	cleaned = strings.Join(fields, " ")

	// 3. 英文统一转为小写
	cleaned = strings.ToLower(cleaned)

	// 4. 移除不影响语义的特殊符号（保留中文常用标点）
	cleaned = specialCharRegex.ReplaceAllString(cleaned, "")

	// 5. 合并连续重复的标点（如！！！→！，，，→，），减少向量噪声
	cleaned = repeatPunctRegex.ReplaceAllString(cleaned, "$1")

	return cleaned
}

func buildWhereFilter(cqr *ChromaQueryRequest) chroma.WhereFilter {
	var conditions []chroma.WhereClause

	// 1. 提供商组织：字符串等值匹配
	if cqr.ProviderOrganization != "" {
		conditions = append(conditions, chroma.EqString("provider_organization", cqr.ProviderOrganization))
	}

	// 2. 能力标志：布尔等值匹配
	if cqr.RequiredCapabilities.Streaming {
		conditions = append(conditions, chroma.EqBool("cap_streaming", true))
	}
	if cqr.RequiredCapabilities.PushNotifications {
		conditions = append(conditions, chroma.EqBool("cap_push_notifications", true))
	}

	// 3. 输入模式：元数据数组包含（数组字段必须用 chroma.K() 包裹键名）
	for _, mode := range cqr.RequiredInputModes {
		conditions = append(conditions, chroma.MetadataContainsString(chroma.K("all_input_modes"), mode))
	}

	// 4. 输出模式：元数据数组包含
	for _, mode := range cqr.RequiredOutputModes {
		conditions = append(conditions, chroma.MetadataContainsString(chroma.K("all_output_modes"), mode))
	}

	// 无任何约束，不启用过滤
	if len(conditions) == 0 {
		return nil
	}

	// 用 AND 组合所有条件，WhereClause 实现了 WhereFilter 接口，可直接返回
	return chroma.And(conditions...)
}

func queryFilter(res chroma.QueryResult, threshold float32) (chroma.QueryResult, error) {
	idGroups := res.GetIDGroups()
	if len(idGroups) == 0 || len(idGroups[0]) == 0 {
		return &chroma.QueryResultImpl{}, nil
	}

	ids := idGroups[0]
	distances := res.GetDistancesGroups()[0]
	metadatas := res.GetMetadatasGroups()[0]
	documents := res.GetDocumentsGroups()[0]

	var filteredIDs chroma.DocumentIDs
	var filteredDistances embeddings.Distances
	var filteredMetadatas chroma.DocumentMetadatas
	var filteredDocuments chroma.Documents

	for i := range ids {
		similarityScore := 1 - float32(distances[i])/2
		if threshold > 0 && similarityScore < threshold {
			continue
		}
		filteredIDs = append(filteredIDs, ids[i])
		filteredDistances = append(filteredDistances, distances[i])
		filteredMetadatas = append(filteredMetadatas, metadatas[i])
		filteredDocuments = append(filteredDocuments, documents[i])
	}

	// 4. 组装返回结果
	filteredResult := &chroma.QueryResultImpl{
		IDLists:         []chroma.DocumentIDs{filteredIDs},
		DistancesLists:  []embeddings.Distances{filteredDistances},
		MetadatasLists:  []chroma.DocumentMetadatas{filteredMetadatas},
		DocumentsLists:  []chroma.Documents{filteredDocuments},
		EmbeddingsLists: res.GetEmbeddingsGroups(),
	}
	return filteredResult, nil
}