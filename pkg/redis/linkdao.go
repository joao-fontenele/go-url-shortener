package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joao-fontenele/go-url-shortener/pkg/common"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

type dao struct {
	conn *redis.Client
}

// NewLinkDao instantiates a dao for Link on redis
func NewLinkDao(conn *redis.Client) shortener.LinkDao {
	return &dao{
		conn: conn,
	}
}

func formatCacheString(slug string) string {
	prefix := common.GetConf().Cache.CachePrefix
	return fmt.Sprintf("%s^l^%s", prefix, slug)
}

func (d *dao) Find(ctx context.Context, slug string) (*shortener.Link, error) {
	link := shortener.Link{}
	str, err := d.conn.Get(ctx, formatCacheString(slug)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, shortener.ErrLinkNotFound
		}
		return nil, err
	}

	json.Unmarshal([]byte(str), &link)

	return &link, nil
}

func (d *dao) Insert(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
	val, _ := json.Marshal(l)
	err := d.conn.Set(
		ctx,
		formatCacheString(l.Slug),
		val,
		time.Duration(common.GetConf().Cache.LinksTTLSeconds)*time.Second,
	).Err()

	return l, err
}

func (d *dao) Update(ctx context.Context, l *shortener.Link) error {
	panic("not yet implemented")
}

func (d *dao) Delete(ctx context.Context, slug string) error {
	key := formatCacheString(slug)

	_, err := d.conn.Del(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}
