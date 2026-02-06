package shorten

import "context"

// ShortenURL generates a short code for the given URL and stores it in the repository.
// It returns the generated code or an error if code generation or storage fails.
// The expire parameter specifies the expiration time in seconds (0 means default expiration).
func (s *shortenURL) ShortenURL(ctx context.Context, url string, expire int) (string, error) {
	code, err := s.keyGen.GenerateCode(urlCodeLength)
	if err != nil {
		return "", err
	}

	ok, err := s.repository.StoreIfNotExists(ctx, code, url, expire)

	switch {
	case err != nil:
		return "", err
	case !ok:
		return "", ErrDuplicatedKey
	}

	return code, nil
}
