package jwt

import "context"

func (s *JwtStorage) DeleteTokenPair(ctx context.Context, access string, refresh string) error {

	// Delete access token
	if err := s.deleteAccess(ctx, access); err != nil {
		return err
	}

	// Delete refresh token
	if err := s.deleteRefresh(ctx, refresh); err != nil {
		return err
	}

	return nil
}

func (s *JwtStorage) deleteAccess(ctx context.Context, access string) error {
	return s.redis.Del(ctx, "acc:"+access).Err()
}

func (s *JwtStorage) deleteRefresh(ctx context.Context, refresh string) error {
	return s.redis.Del(ctx, "ref:"+refresh).Err()
}
