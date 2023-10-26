package service

import (
	"github.com/prclin/alumni-circle/dao"
	_error "github.com/prclin/alumni-circle/error"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"hash/fnv"
	"math"
	"net/http"
	"sort"
	"strconv"
)

func GetBreakFeed(accountId uint64, latestTime int64, count int) ([]model.Break, error) {
	//获取用户标签
	tagDao := dao.NewTagDao(global.Datasource)
	accountTags, err := tagDao.SelectEnabledAccountTagByAccountId(accountId)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.NewWithCode(http.StatusInternalServerError)
	}
	//随机获取发布的帖子
	breakDao := dao.NewBreakDao(global.Datasource)
	breakIds, err := breakDao.SelectApprovedIdsRandomlyBefore(latestTime, accountId, 100*count)
	filteredBreakIds := make([]uint64, 0, len(breakIds))
	//去除已推荐的帖子，通过redis实现的布隆过滤器
	for _, value := range breakIds {
		recommended := isRecommendedTo(accountId, value)
		if recommended {
			continue
		}
		filteredBreakIds = append(filteredBreakIds, value)
	}
	breakIds = filteredBreakIds
	//获取帖子标签
	exactCount := len(breakIds)
	breakTagMap := make(map[uint64][]model.TTag, exactCount)
	for _, value := range breakIds {
		tag, err := tagDao.SelectEnabledBreakTagByBreakId(value)
		if err != nil {
			global.Logger.Debug(err)
			return nil, _error.NewWithCode(http.StatusInternalServerError)
		}
		breakTagMap[value] = tag
	}

	//计算相似度向量
	similarityMap := tagCosineSimilarityMap(accountTags, breakTagMap)

	//相似度排序获取top count
	breakSimilarities := make(BreakSimilaritySlice, 0, exactCount)
	for key, value := range similarityMap {
		similarity := BreakSimilarity{BreakId: key, Similarity: value}
		breakSimilarities = append(breakSimilarities, similarity)
		breakSimilarities = append(breakSimilarities, similarity)
	}
	sort.Sort(breakSimilarities)
	count = int(math.Min(float64(count), float64(exactCount)))
	topCount := make([]uint64, 0, count)
	for _, value := range breakSimilarities[:count] {
		topCount = append(topCount, value.BreakId)
	}
	//将top count的break id存入redis布隆过滤器
	err = recordRecommendedBreaks(accountId, topCount)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.NewWithCode(http.StatusInternalServerError)
	}
	//查询top count课间并返回
	tBreaks, err := breakDao.SelectByIds(topCount)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.NewWithCode(http.StatusInternalServerError)
	}
	breaks := make([]model.Break, 0, len(tBreaks))
	shotDao := dao.NewShotDao(global.Datasource)
	for _, value := range breaks {
		shots, err := shotDao.SelectShotsByBreakId(value.Id)
		if err != nil {
			global.Logger.Debug(err)
			return nil, _error.NewWithCode(http.StatusInternalServerError)
		}
		value.Shots = shots
		value.Tags = breakTagMap[value.Id]
	}

	return breaks, nil
}

type RecommendBloomFilter struct {
	Key         string
	bitSize     uint64
	hashFuncNum int
}

// Add 向过滤器中添加课间id
func (filter RecommendBloomFilter) Add(breakId uint64) error {
	for i := 0; i < filter.hashFuncNum; i++ {
		offset := int64(filter.hash(strconv.FormatUint(breakId, 10), i) % filter.bitSize)
		_, err := dao.SetBit(filter.Key, offset, 1)
		if err != nil {
			return err
		}
	}
	return nil
}

// Contains 过滤器中是否有某个课间id
func (filter RecommendBloomFilter) Contains(breakId uint64) (bool, error) {
	for i := 0; i < filter.hashFuncNum; i++ {
		offset := int64(filter.hash(strconv.FormatUint(breakId, 10), i) % filter.bitSize)
		value, err := dao.GetBit(filter.Key, offset)
		if err != nil {
			return false, err
		}
		if value == 0 {
			return false, nil
		}
	}
	return true, nil
}
func (filter RecommendBloomFilter) hash(item string, seed int) uint64 {
	h := fnv.New64a()
	h.Write([]byte(item))
	hash := h.Sum64()
	return hash + uint64(seed)
}

func NewRecommendBloomFilter(key string) *RecommendBloomFilter {
	return &RecommendBloomFilter{Key: key, bitSize: 5000, hashFuncNum: 3}
}

// recordRecommendedBreaks 记录已推荐给指定用户的课间id
func recordRecommendedBreaks(accountId uint64, breakIds []uint64) error {
	filter := NewRecommendBloomFilter("break:recommended:" + strconv.FormatUint(accountId, 10))
	for _, id := range breakIds {
		err := filter.Add(id)
		if err != nil {
			return err
		}
	}
	return nil
}

// isRecommendedTo 判断一个课间是否被推荐给指定用户
func isRecommendedTo(accountId, breakId uint64) bool {
	filter := NewRecommendBloomFilter("break:recommended:" + strconv.FormatUint(accountId, 10))
	contains, err := filter.Contains(breakId)
	if err != nil {
		return false
	}
	return contains
}

// 确保BreakSimilarity实现sort.Interface
var _ sort.Interface = BreakSimilaritySlice{}

type BreakSimilarity struct {
	BreakId    uint64
	Similarity float64
}

type BreakSimilaritySlice []BreakSimilarity

func (bss BreakSimilaritySlice) Len() int {
	return len(bss)
}

func (bss BreakSimilaritySlice) Less(i, j int) bool {
	return bss[i].Similarity > bss[j].Similarity //降序
}

func (bss BreakSimilaritySlice) Swap(i, j int) {
	bss[i], bss[j] = bss[j], bss[i]
}

// tagCosineSimilarityMap 计算账户tag与课间tag集合的余弦相似度集合
func tagCosineSimilarityMap(accountTags []model.TTag, tagMap map[uint64][]model.TTag) map[uint64]float64 {
	similarityMap := make(map[uint64]float64, len(tagMap))
	//逐个计算相似度
	for key, value := range tagMap {
		similarity := tagCosineSimilarity(accountTags, value)
		similarityMap[key] = similarity
	}
	return similarityMap
}

// tagCosineSimilarity 计算两个向量的余弦相似度
func tagCosineSimilarity(tags1, tags2 []model.TTag) float64 {
	//构造并集
	tagsSet := make(map[uint32]struct{}, len(tags1)+len(tags2))
	for _, tag := range tags1 {
		tagsSet[tag.Id] = struct{}{}
	}
	for _, tag := range tags2 {
		tagsSet[tag.Id] = struct{}{}
	}

	//创建标签向量
	v1 := make([]float64, len(tagsSet))
	v2 := make([]float64, len(tagsSet))
	// 填充标签向量
	i := 0
	for tagId, _ := range tagsSet {
		//v1
		if containsTag(tags1, tagId) {
			v1[i] = 1
		} else {
			v1[i] = 0
		}
		//v2
		if containsTag(tags2, tagId) {
			v2[i] = 1
		} else {
			v2[i] = 0
		}
		i++
	}

	// 计算余弦相似度
	dotProduct := 0.0 //点积
	v1m := 0.0        //v1模
	v2m := 0.0        //v2模
	tsl := len(tagsSet)
	for j := 0; j < tsl; j++ {
		dotProduct += v1[j] * v2[j]
		v1m += v1[j] * v1[j]
		v2m += v2[j] * v2[j]
	}

	v1m = math.Sqrt(v1m)
	v2m = math.Sqrt(v2m)

	return dotProduct / (v1m * v2m)
}

func containsTag(tags []model.TTag, tagId uint32) bool {
	for _, value := range tags {
		if value.Id == tagId {
			return true
		}
	}
	return false
}

func DeleteBreak(tBreak model.TBreak) error {
	breakDao := dao.NewBreakDao(global.Datasource)
	return breakDao.DeleteByIdAndAccountId(tBreak.Id, tBreak.AccountId)
}

func UpdateBreakVisibility(tBreak model.TBreak) error {
	bd := dao.NewBreakDao(global.Datasource)
	return bd.UpdateVisibilityBy(tBreak)
}

func PublishBreak(tBreak model.TBreak, shotIds, topicIds []uint64) (model.Break, error) {
	tx := global.Datasource.Begin()
	defer tx.Commit()
	bd := dao.NewBreakDao(tx)
	//创建课间
	breakId, err := bd.InsertBy(tBreak)
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	//绑定图片
	shotBindings := make([]model.TShotBinding, 0, len(shotIds))
	for index, shotId := range shotIds {
		shotBindings = append(shotBindings, model.TShotBinding{BreakId: breakId, ImageId: shotId, Order: uint8(index)})
	}
	sd := dao.NewShotDao(tx)
	err = sd.BatchInsertBy(shotBindings)
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	//绑定话题
	topicBindings := make([]model.TTopicBinding, 0, len(topicIds))
	for _, topicId := range topicIds {
		topicBindings = append(topicBindings, model.TTopicBinding{BreakId: breakId, TopicId: topicId})
	}
	td := dao.NewTopicDao(tx)
	err = td.BatchInsertBindingBy(topicBindings)
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	//获取课间
	var _break model.Break
	tb, err := bd.SelectById(breakId) //基本信息
	_break.TBreak = tb
	shots, err := sd.SelectShotsByBreakId(breakId) //镜头
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	_break.Shots = shots
	topics, err := td.SelectTopicsByBreakId(breakId) //话题
	if err != nil {
		tx.Rollback()
		return model.Break{}, err
	}
	_break.Topics = topics

	return _break, nil
}
