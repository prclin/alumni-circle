package service

import (
	"context"
	"github.com/prclin/alumni-circle/dao"
	_error "github.com/prclin/alumni-circle/error"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
	"github.com/prclin/alumni-circle/util"
	"github.com/redis/go-redis/v9"
	"hash/fnv"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// AcquireLikedBreak 获取用户点赞
func AcquireLikedBreak(acquirer, acquiree uint64, pagination model.Pagination) ([]model.Break, error) {
	breakDao := dao.NewBreakDao(global.Datasource)
	//已点赞课间id
	breakIds, err := breakDao.SelectLikedIdsBy(acquiree, pagination)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.InternalServerError
	}

	//缓存中已点赞课程id
	cachedMap, err := getCachedLikes(acquiree)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.InternalServerError
	}

	//获取交集，并且以缓存中点赞为准
	for _, breakId := range breakIds {
		_, ok := cachedMap[breakId]
		if ok {
			continue
		}
		cachedMap[breakId] = 1
	}

	likedBreakIds := make([]uint64, 0, len(cachedMap))
	for key, value := range cachedMap {
		if value == 1 {
			likedBreakIds = append(likedBreakIds, key)
		}
	}

	//查询已点赞课间
	tBreaks, err := breakDao.SelectByIds(likedBreakIds)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.InternalServerError
	}

	//映射为结果
	breaks := make([]model.Break, 0, len(tBreaks))
	shotDao := dao.NewShotDao(global.Datasource)
	tagDao := dao.NewTagDao(global.Datasource)
	for _, tBreak := range tBreaks {
		//不可见则跳过
		visibility, err1 := getBreakVisibility(acquirer, acquiree)
		if err1 != nil {
			global.Logger.Warn(err1)
		}
		if tBreak.Visibility < visibility {
			continue
		}
		shots, err1 := shotDao.SelectShotsByBreakId(tBreak.Id) //镜头
		if err1 != nil {
			global.Logger.Warn(err1)
		}
		tags, err1 := tagDao.SelectEnabledByBreakId(tBreak.Id) //标签
		if err1 != nil {
			global.Logger.Warn(err1)
		}
		info, err1 := GetAccountInfo(acquirer, acquiree) //账户信息
		if err1 != nil {
			global.Logger.Warn(err1)
		}
		//点赞数
		growth, _ := dao.HGet("break_like_growth", strconv.FormatUint(tBreak.Id, 10))
		tBreak.LikeCount = tBreak.LikeCount + uint32(util.IgnoreError(strconv.ParseUint(growth, 10, 32)))
		breaks = append(breaks, model.Break{TBreak: tBreak, Shots: shots, Tags: tags, AccountInfo: info, Liked: true})
	}

	return breaks, nil
}

func getCachedLikes(accountId uint64) (map[uint64]uint8, error) {
	hashMap, err := dao.HScan("break_likes", strconv.FormatUint(accountId, 10)+":*")
	if err != nil {
		return nil, err
	}
	likes := make(map[uint64]uint8, len(hashMap))
	for key, value := range hashMap {
		likes[util.IgnoreError(strconv.ParseUint(strings.Split(key, ":")[1], 10, 64))] = uint8(util.IgnoreError(strconv.ParseUint(value, 10, 8)))
	}
	return likes, nil
}

// AcquireBreakList 获取账户课间列表
func AcquireBreakList(acquirer, acquiree uint64, pagination model.Pagination) ([]model.Break, error) {
	//获取账户信息
	info, err := GetAccountInfo(acquirer, acquiree)
	if err != nil {
		return nil, _error.InternalServerError
	}

	//获取可见范围
	visibility, err := getBreakVisibility(acquirer, acquiree)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.InternalServerError
	}

	//获取课间
	breakDao := dao.NewBreakDao(global.Datasource)
	tBreaks, err := breakDao.SelectByAccountIdAndVisibility(acquiree, visibility, pagination)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.InternalServerError
	}

	breaks := make([]model.Break, 0, len(tBreaks))

	shotDao := dao.NewShotDao(global.Datasource)
	tagDao := dao.NewTagDao(global.Datasource)
	for _, tBreak := range tBreaks {
		shots, err1 := shotDao.SelectShotsByBreakId(tBreak.Id) //镜头
		if err1 != nil {
			global.Logger.Warn(err1)
		}
		tags, err1 := tagDao.SelectEnabledByBreakId(tBreak.Id) //标签
		if err1 != nil {
			global.Logger.Warn(err1)
		}
		//是否点赞
		var liked bool
		action, err1 := dao.HGet("break_likes", strconv.FormatUint(acquirer, 10)+":"+strconv.FormatUint(tBreak.Id, 10))
		if err1 != nil && err != redis.Nil {
			global.Logger.Warn(err1)
		}
		liked = action == "1"
		if err == redis.Nil {
			liked = breakDao.IsLiked(acquirer, tBreak.Id)
		}
		breaks = append(breaks, model.Break{TBreak: tBreak, Shots: shots, Tags: tags, AccountInfo: info, Liked: liked})
	}
	return breaks, nil
}

// getBreakVisibility 获取用户对于指定用户的课间可见范围
func getBreakVisibility(acquirer uint64, acquiree uint64) (uint8, error) {
	if acquirer == acquiree {
		return 0, nil
	}
	//获取关系，以推断可见性
	followDao := dao.NewFollowDao(global.Datasource)

	var visibility uint8 = 3 //所有人可见

	followed, err := followDao.IsFollowed(acquirer, acquiree)
	if err != nil {
		return visibility, err
	}

	if followed { //粉丝可见
		visibility--

		beFollowed, err := followDao.IsFollowed(acquiree, acquirer)
		if err != nil {
			return visibility, nil
		}

		if beFollowed { //互关可见
			visibility--
		}
	}

	return visibility, err
}

// FlushBreakLikes 将点赞落库
func FlushBreakLikes() {
	const likesKey = "break_likes"
	const likeGrowthKey = "break_like_growth"
	expireTime := time.Hour //重置过期时间
	defer func() {
		//重新计时
		err := dao.SetString("expired_"+likesKey, expireTime.String(), expireTime)
		if err != nil {
			global.Logger.Error("无法重新倒计时，请及时处理")
		}
	}()

	//获取课间点赞
	likeMap, err := dao.HGetAll(likesKey)
	if err != nil && err != redis.Nil {
		global.Logger.Debug(err)
		//设置 expired_break_likes 5分钟后过期
		expireTime = 5 * time.Minute
		return
	}
	//获取课间点赞数
	likeGrowthMap, err := dao.HGetAll(likeGrowthKey)
	if err != nil && err != redis.Nil {
		global.Logger.Debug(err)
		//设置 expired_break_likes 5分钟后过期
		expireTime = 5 * time.Minute
		return
	}

	//点赞落库
	likes := make([]model.TBreakLike, 0, len(likeMap)/2)   //点赞
	unlikes := make([]model.TBreakLike, 0, len(likeMap)/2) //取消点赞
	for key, value := range likeMap {
		split := strings.Split(key, ":")
		accountId := util.IgnoreError(strconv.ParseUint(split[0], 10, 64))
		breakId := util.IgnoreError(strconv.ParseUint(split[1], 10, 64))
		breakLike := model.TBreakLike{AccountId: accountId, BreakId: breakId}
		switch value {
		case "0":
			unlikes = append(unlikes, breakLike)
			break
		case "1":
			likes = append(likes, breakLike)
			break
		}
	}
	tx := global.Datasource.Begin() //开启事务
	defer tx.Commit()
	breakDao := dao.NewBreakDao(tx)
	err = breakDao.BatchInsertLikeBy(likes) //点赞
	if err != nil {                         //重置过期时间，一会儿再试
		global.Logger.Debug(err)
		tx.Rollback()
		return
	}
	err = breakDao.BatchDeleteLikeBy(unlikes) //取消点赞
	if err != nil {
		global.Logger.Debug(err)
		tx.Rollback()
		return
	}

	//点赞数落库
	increases := make(map[uint64]int, len(likeGrowthMap))
	for key, value := range likeGrowthMap {
		increases[util.IgnoreError(strconv.ParseUint(key, 10, 64))] = int(util.IgnoreError(strconv.ParseInt(value, 10, 64)))
	}
	err = breakDao.BatchIncreaseLikeCount(increases)
	if err != nil {
		global.Logger.Debug(err)
		tx.Rollback()
		return
	}

	//删除过时
	_, err1 := dao.DeleteKey(likesKey)
	if err1 != nil {
		global.Logger.Warn("无法删除过时点赞，请及时处理")
	}
	_, err2 := dao.DeleteKey(likeGrowthKey)
	if err2 != nil {
		global.Logger.Warn("无法删除过时点赞数，请及时处理")
	}
}

const likeScript = `
	local lockKey=KEYS[1]
	--获取锁
	local lockResult = redis.pcall("GET",lockKey)
	if type(lockResult) == 'table' and lockResult.err then --获取锁发生错误，打印调试信息，并返回错误信息
  		redis.log(redis.LOG_NOTICE, "get break_lock failed", lockResult.err)
		return {err = "获取锁时发生错误，注意lock key必须是一个string类型的key,具体错误为：" + lockResult.err}
	end
	--没获取到锁，直接返回
	if not lockResult then
		return false
	end
	
	--点赞或取消
	local likeResult = redis.pcall("HSET",KEYS[2],ARGV[1],ARGV[2])
	if type(likeResult) == 'table' and likeResult.err then --点赞或取消发生错误，打印调试信息，并返回错误信息
	  	redis.log(redis.LOG_NOTICE, "set break_like field failed", likeResult.err)
		return {err = "点赞或取消时发生错误，注意break_like key必须是一个hash类型的key,具体错误为：" + likeResult.err}
	end

	--点赞增长数
	local growthResult = redis.pcall("HINCRBY",KEYS[3],ARGV[3],ARGV[4])
	if type(growthResult) == 'table' and growthResult.err then --更新点赞增长数发生错误，打印调试信息，并返回错误信息
		redis.log(redis.LOG_NOTICE, "set break_like_growth field failed", growthResult.err)
		return {err = "更新点赞增长数时发生错误，注意break_like_growth key必须是一个hash类型的key,具体错误为：" + growthResult.err}
	end
	
	return true
	`

// LikeBreak 点赞或取消点赞课间
//
// 使用lua脚本，获取锁并写到redis
func LikeBreak(bl *model.TBreakLike, action int) error {
	const lockKey = "expired_break_likes"
	const likesKey = "break_likes"
	const likeGrowthKey = "break_like_growth"
	//构建脚本
	script := redis.NewScript(likeScript)
	//有10次重试机会
	var err error
	for i := 0; i < 10; i++ {
		err = script.Run(context.Background(), global.RedisClient, []string{lockKey, likesKey, likeGrowthKey}, bl.String(), action, strconv.FormatUint(bl.BreakId, 10), util.Ternary(action == 1, 1, -1)).Err()
		if err != nil {
			if err == redis.Nil && i < 9 { //抢锁失败，休眠50ms，最后一次不休眠
				time.Sleep(50 * time.Millisecond)
			} else { //其它错误直接结束重试
				break
			}
		} else { //成功直接退出循环
			break
		}
	}

	if err != nil {
		global.Logger.Debug(err)
		return util.Ternary(err == redis.Nil, _error.InternalServerError, _error.NewServerError("服务器忙，请稍后重试"))
	}

	return nil
}

func GetBreakFeed(accountId uint64, latestTime int64, count int) ([]model.Break, error) {
	//获取用户标签
	tagDao := dao.NewTagDao(global.Datasource)
	accountTags, err := tagDao.SelectEnabledByAccountId(accountId)
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
		tag, err := tagDao.SelectEnabledByBreakId(value)
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
		info, err := GetAccountInfo(accountId, value.AccountId)
		if err != nil {
			global.Logger.Debug(err)
			return nil, _error.NewWithCode(http.StatusInternalServerError)
		}

		value.Shots = shots
		value.AccountInfo = info
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
	for tagId := range tagsSet {
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

func PublishBreak(tBreak model.TBreak, shotIds []uint64, tagIds []uint32) (*model.Break, error) {
	tx := global.Datasource.Begin()
	defer tx.Commit()
	bd := dao.NewBreakDao(tx)
	//创建课间
	breakId, err := bd.InsertBy(tBreak)
	if err != nil {
		tx.Rollback()
		return nil, err
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
		return nil, err
	}
	//绑定标签
	tagBindings := make([]model.TBreakTagBinding, 0, len(tagIds))
	for _, tagId := range tagIds {
		tagBindings = append(tagBindings, model.TBreakTagBinding{BreakId: breakId, TagId: tagId})
	}
	td := dao.NewTagDao(tx)
	err = td.BatchInsertBreakTagBindingBy(tagBindings)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//获取课间
	_break := &model.Break{}
	tb, err := bd.SelectById(breakId) //基本信息
	_break.TBreak = *tb
	shots, err := sd.SelectShotsByBreakId(breakId) //镜头
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_break.Shots = shots
	tags, err := td.SelectEnabledByBreakId(breakId) //话题
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_break.Tags = tags

	return _break, nil
}
