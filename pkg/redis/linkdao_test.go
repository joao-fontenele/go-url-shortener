package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/joao-fontenele/go-url-shortener/pkg/common"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

func testMain(m *testing.M) int {
	var err error

	os.Chdir("/usr/src/app")

	err = common.LoadConfs()
	if err != nil {
		fmt.Println("failed to load configs: %w", err)
		return 1
	}

	env := common.GetConf().Env
	if env != "test" {
		fmt.Println("don't run these tests on non dev environment")
		return 1
	}

	closeDb, err := Connect()
	if err != nil {
		return 1
	}

	defer closeDb()

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func seedDB(conn *redis.Client) error {
	key := formatCacheString("a1CDz")
	err := conn.Set(
		context.Background(),
		key,
		`{"slug":"a1CDz","url":"https://www.google.com","createdAt":"2020-05-01T00:00:00.000Z"}`,
		1*time.Duration(time.Minute),
	).Err()

	return err
}

func truncateDB(conn *redis.Client) error {
	var cursor uint64
	cachePrefix := fmt.Sprintf("%s*", common.GetConf().Cache.CachePrefix)
	ctx := context.Background()

	for {
		var keys []string
		var err error
		keys, cursor, err = conn.Scan(ctx, cursor, cachePrefix, 10).Result()

		if err != nil {
			return err
		}

		if len(keys) == 0 {
			return nil
		}

		_, err = conn.Del(ctx, keys...).Result()

		if err != nil {
			return err
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func TestFind(t *testing.T) {
	conn := GetConnection()
	err := truncateDB(conn)
	if err != nil {
		t.Fatalf("error truncating test database: %v", err)
	}

	err = seedDB(conn)
	if err != nil {
		t.Fatalf("error seeding database: %v", err)
	}

	tt := []struct {
		Name string
		Slug string
		Want *shortener.Link
	}{
		{
			Name: "FoundSlug",
			Slug: "a1CDz",
			Want: &shortener.Link{
				URL:       "https://www.google.com",
				Slug:      "a1CDz",
				CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Name: "NotFoundSlug",
			Slug: "niull",
			Want: nil,
		},
	}

	dao := NewLinkDao(conn)

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			got, err := dao.Find(context.Background(), test.Slug)

			if err != nil {
				t.Fatalf("failed to find requested link: %v", err)
			}

			if diff := cmp.Diff(test.Want, got); diff != "" {
				t.Errorf("failed to fetch expected link (-want +got):\n%s", diff)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	conn := GetConnection()

	inserted := &shortener.Link{
		Slug:      "f00Bar",
		URL:       "https://wwww.duckduckgo.com",
		CreatedAt: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
	}
	key := formatCacheString(inserted.Slug)
	ctx := context.Background()

	dao := NewLinkDao(conn)
	_, err := dao.Insert(ctx, inserted)

	if err != nil {
		t.Fatalf("failed to insert new link: %v", err)
	}

	rawGot, err := conn.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("failed to query database for inserted link: %v", err)
	}
	var got shortener.Link
	json.Unmarshal([]byte(rawGot), &got)

	if diff := cmp.Diff(inserted, &got); diff != "" {
		t.Errorf("failed to insert link correctly (-want +got):\n%s", diff)
	}

	currTTL, err := conn.TTL(ctx, key).Result()
	if err != nil {
		t.Errorf("failed to query for inserted key ttl: %v", err)
	}

	if currTTL == -1 {
		t.Error("should have set key ttl, but it has no ttl")
	}

	expectedDur := time.Duration(common.GetConf().Cache.LinksTTLSeconds) * time.Second
	if expectedDur-currTTL > 1 {
		t.Errorf("ttl is set to the wrong value, (want, got): (%v, %v)", expectedDur, currTTL)
	}
}

func TestDelete(t *testing.T) {
	conn := GetConnection()

	dao := NewLinkDao(conn)

	tt := []struct {
		Name string
		Key  string
	}{
		{
			Name: "DeleteFound",
			Key:  "a1CDz",
		},
		{
			Name: "DeleteNotFound",
			Key:  "h3ll0",
		},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			ctx := context.Background()

			err := dao.Delete(ctx, test.Key)
			if err != nil {
				t.Fatalf("failed to delete slug: %v", err)
			}

			_, err = conn.Get(ctx, formatCacheString(test.Key)).Result()

			if err == nil {
				t.Fatal("key should have been deleted")
			}

			if !errors.Is(err, redis.Nil) {
				t.Fatalf("unexpected error while getting cached key: %v", err)
			}
		})
	}
}
