package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/panzerhomer/banner/internal/config"
	"github.com/panzerhomer/banner/internal/domain"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Redis struct {
	client *redis.Client
}

func New(cfg *config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("redis not send PONG ("+pong+"), err: ", err)
	}

	return &Redis{client: client}, nil
}

func (r *Redis) SaveBanner(bannerID int64, banner domain.Banner) error {
	ttl := time.Minute * 5
	tagIds := makeString(banner.TagIds)
	key := tagIds + fmt.Sprint(banner.FeatureID)

	content, _ := json.Marshal(banner)
	// fmt.Println("content", string(content))
	if err := r.client.Set(ctx, key, string(content), ttl).Err(); err != nil {
		return errors.Wrapf(err, "error when try save in cache with key: %s", key)
	}

	return nil
}

func (r *Redis) GetBanner(tagIDs []int64, featureID int64) (*domain.Banner, error) {
	tagIds := makeString(tagIDs)
	key := tagIds + fmt.Sprint(featureID)

	var content string
	if err := r.client.Get(ctx, key).Scan(&content); err != nil {
		return nil, errors.Wrapf(err, "error when try get cache with key: %s", key)
	}

	// fmt.Println("GetBanner", content, tagIDs, key)

	var bannerCached domain.Banner
	err := json.Unmarshal([]byte(content), &bannerCached)
	if err != nil {
		return nil, errors.New("banner cached unmarshaling failed")
	}

	return &bannerCached, nil
}

func makeString(nums []int64) string {
	var strNumbers []string
	for _, num := range nums {
		strNumbers = append(strNumbers, strconv.Itoa(int(num)))
	}

	result := strings.Join(strNumbers, ",")

	return result
}
