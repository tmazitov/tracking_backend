package jwt

import (
	"context"
	"time"
)

var (
	accessLifeTime  = time.Minute * 15
	refreshLifeTime = time.Hour * 24 * 30
)

func (j *JwtStorage) accessToArchive(ctx context.Context, access string) error {
	return j.redis.Set(ctx, "acc:"+access, 1, accessLifeTime).Err()
}

func (j *JwtStorage) refreshToArchive(ctx context.Context, refresh string) error {
	return j.redis.Set(ctx, "ref:"+refresh, 1, refreshLifeTime).Err()
}

func (j *JwtStorage) recordToArchive(ctx context.Context, tokens TokenPair) error {

	// record access token to cache
	if err := j.accessToArchive(ctx, tokens.Access); err != nil {
		return err
	}

	// record refresh token to cache
	if err := j.refreshToArchive(ctx, tokens.Refresh); err != nil {
		return err
	}

	return nil
}
