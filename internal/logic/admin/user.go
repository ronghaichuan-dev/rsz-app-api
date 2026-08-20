package admin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/admin/v1"
)

func (s *sAdmin) Login(ctx context.Context, in v1.SystemUserLoginInput) (*v1.SystemUserLoginOutput, error) {
	token, err := generateAdminToken(1)
	if err != nil {
		return nil, err
	}
	return &v1.SystemUserLoginOutput{
		UserId: 1,
		Token:  token,
		Name:   in.Username,
	}, nil
}

func generateAdminToken(userId uint64) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", gerror.Wrap(err, "generate admin token failed")
	}
	return fmt.Sprintf("admin_%d_%d_%s", userId, time.Now().Unix(), hex.EncodeToString(buf)), nil
}
