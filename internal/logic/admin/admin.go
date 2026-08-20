package admin

import "rslytics-app-api/internal/service"

type sAdmin struct {
}

func init() {
	service.RegisterAdmin(NewAdmin())
}

func NewAdmin() *sAdmin {
	return &sAdmin{}
}
