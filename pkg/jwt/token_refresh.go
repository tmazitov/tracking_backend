package jwt

import "context"

func (s *JwtStorage) RefreshIsExists(ctx context.Context, token string) error {
	return s.isExists(ctx, "ref:", token)

}

func (s *JwtStorage) ValidateRefresh(ctx context.Context, token string) (*AccessClaims, error) {

	// Check if token is exists
	if err := s.RefreshIsExists(ctx, token); err != nil {
		return nil, err
	}

	return s.verifyToken(ctx, token)
}
