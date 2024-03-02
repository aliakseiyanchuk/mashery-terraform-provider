package mashres

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/redis/go-redis/v9"
	"io"
	"time"
)

var redisClient *redis.Client
var DefaultCacheDuration = time.Hour

type CacheStruct[TIdent, TValue any] struct {
	Ident TIdent
	Value TValue
}

func InitRedisClient(url string) error {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return err
	}

	redisClient = redis.NewClient(opts)
	return nil
}

func GetFromCache[TIdent, TValue any](ctx context.Context, key string) (*TIdent, *TValue) {
	if redisClient == nil {
		return nil, nil
	}

	val, readErr := redisClient.Get(ctx, key).Result()
	if readErr != nil {
		tflog.Error(ctx, "reading redis cache failed", map[string]interface{}{"key": key, "error": readErr})
		return nil, nil
	}

	if encBytes, base64Err := base64.StdEncoding.DecodeString(val); base64Err != nil {
		tflog.Error(ctx, "decoding base64 failed", map[string]interface{}{"key": key, "error": base64Err})
		return nil, nil
	} else {
		gzStream, gzipErr := gzip.NewReader(bytes.NewReader(encBytes))
		if gzipErr != nil {
			tflog.Error(ctx, "decoding gzip failed", map[string]interface{}{"key": key, "error": gzipErr})
			return nil, nil
		}

		jsonBytes, jsonReadErr := io.ReadAll(gzStream)
		if jsonReadErr != nil {
			tflog.Error(ctx, "reading gzip failed", map[string]interface{}{"key": key, "error": jsonReadErr})
			return nil, nil
		}

		var rv CacheStruct[TIdent, TValue]
		if jsonErr := json.Unmarshal(jsonBytes, &rv); jsonErr != nil {
			tflog.Error(ctx, "decoding json failed", map[string]interface{}{"key": key, "error": jsonReadErr})
			return nil, nil
		}

		return &rv.Ident, &rv.Value
	}
}

func StoreInCacheDefault[TIdent, TValue any](ctx context.Context, key string, ident TIdent, value TValue) {
	StoreInCache(ctx, key, ident, value, DefaultCacheDuration)
}

func StoreInCache[TIdent, TValue any](ctx context.Context, key string, ident TIdent, value TValue, dur time.Duration) {
	if redisClient == nil {
		return
	}

	var rv CacheStruct[TIdent, TValue]
	rv.Ident = ident
	rv.Value = value

	jsonBytes, _ := json.Marshal(&rv)

	var gzipBytes bytes.Buffer
	w := gzip.NewWriter(&gzipBytes)
	w.Write(jsonBytes)
	w.Close()

	base64Out := base64.StdEncoding.EncodeToString(gzipBytes.Bytes())
	redisClient.Set(ctx, key, base64Out, dur)
}
