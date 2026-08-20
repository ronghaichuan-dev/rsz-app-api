package kids

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	v1 "rslytics-app-api/internal/api/kids/v1"
	"rslytics-app-api/internal/consts"
	"rslytics-app-api/internal/utils"
)

// ListNotifications 从数据库查询通知列表，支持按成员和未读状态筛选。
func (s *sKids) ListNotifications(ctx context.Context, in v1.NotificationListInput) (*v1.NotificationListOutput, error) {
	model := utils.KidsDB(ctx).Model(consts.KidsNotificationTable).Ctx(ctx)
	if in.MemberId > 0 {
		model = model.Where("member_id IN(?)", []uint64{0, in.MemberId})
	}
	if in.Unread {
		model = model.Where("is_read", 0)
	}
	records, err := model.OrderDesc("id").All()
	if err != nil {
		return nil, err
	}
	out := &v1.NotificationListOutput{}
	for _, record := range records {
		out.List = append(out.List, notificationFromDB(record))
	}
	return out, nil
}

// ReadNotification 将指定通知持久化标记为已读。
func (s *sKids) ReadNotification(ctx context.Context, in v1.NotificationReadInput) (*v1.NotificationReadOutput, error) {
	db := utils.KidsDB(ctx)
	result, err := db.Model(consts.KidsNotificationTable).Ctx(ctx).Where("id", in.Id).Data(map[string]any{"is_read": 1}).Update()
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "notification not found")
	}
	record, err := db.Model(consts.KidsNotificationTable).Ctx(ctx).Where("id", in.Id).One()
	if err != nil {
		return nil, err
	}
	return &v1.NotificationReadOutput{Notification: notificationFromDB(record)}, nil
}

// createNotification 持久化创建一条通知，并返回创建后的通知结构。
func createNotification(ctx context.Context, tx gdb.TX, memberId uint64, notificationType, title, content string) (v1.NotificationItem, error) {
	id, err := tx.Model(consts.KidsNotificationTable).Ctx(ctx).Data(map[string]any{
		"member_id":         memberId,
		"notification_type": notificationType,
		"title":             title,
		"content":           content,
		"is_read":           0,
	}).InsertAndGetId()
	if err != nil {
		return v1.NotificationItem{}, err
	}
	record, err := tx.Model(consts.KidsNotificationTable).Ctx(ctx).Where("id", id).One()
	if err != nil {
		return v1.NotificationItem{}, err
	}
	return notificationFromDB(record), nil
}

// notificationFromDB 将数据库通知记录转换为接口响应结构。
func notificationFromDB(record gdb.Record) v1.NotificationItem {
	return v1.NotificationItem{
		Id:        record["id"].Uint64(),
		MemberId:  record["member_id"].Uint64(),
		Type:      record["notification_type"].String(),
		Title:     record["title"].String(),
		Content:   record["content"].String(),
		Read:      record["is_read"].Int() == 1,
		CreatedAt: utils.ParseDBTime(record["created_at"].Val()),
	}
}
