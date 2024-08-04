package cache

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"plugin_onenet/model"
	"time"
)

var REDIS *redis.Client

func RedisInit() {
	REDIS = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),     // Redis 服务器地址
		Password: viper.GetString("redis.password"), // 没有密码时保持为空
		DB:       viper.GetInt("redis.db"),          // 使用默认的 DB
	})
}

// SetDeviceInfo
// @description 设置待连接的设备信息
func SetDeviceInfo(ctx context.Context, productId, deviceName string) error {
	productCacheKey := fmt.Sprintf(viper.GetString("onenet.product_cache_key"), productId)
	luaScript := `
    local exists = redis.call('ZSCORE', KEYS[1], ARGV[1])
    if not exists then
        redis.call('ZADD', KEYS[1], ARGV[2], ARGV[1])
    end
    return exists
    `
	_, err := REDIS.Eval(ctx, luaScript, []string{productCacheKey}, deviceName, time.Now().Unix()).Result()
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// GetDeviceList
// @description 获取待添加设备列表
func GetDeviceList(ctx context.Context, productId string, page, pageSize int64) (int64, []model.DeviceItem, error) {
	productCacheKey := fmt.Sprintf(viper.GetString("onenet.product_cache_key"), productId)
	var (
		total int64
		list  []model.DeviceItem
		err   error
	)
	total, err = REDIS.ZCard(ctx, productCacheKey).Result()
	if err != nil {
		logrus.Error(err)
		return total, list, err
	}
	start := (page - 1) * pageSize
	end := start + pageSize - 1
	deviceNumbers, err := REDIS.ZRangeWithScores(ctx, productCacheKey, start, end).Result()
	if err != nil {
		logrus.Error(err)
		return total, list, err
	}
	for _, v := range deviceNumbers {
		list = append(list, model.DeviceItem{
			DeviceNumber: fmt.Sprintf(viper.GetString("onenet.device_number_key"), productId, v.Member.(string)),
			Description:  fmt.Sprintf("%f", v.Score),
			DeviceName:   v.Member.(string),
		})
	}
	return total, list, nil
}
