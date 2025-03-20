package usecase

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"github.com/qbhy/goal-piplin/app/models"
)

func DeleteGroups(id any) error {
	_, err := models.UserGroups().WhereIn("group_id", id).DeleteE()
	if err != nil {
		return err
	}

	_, err = models.Groups().WhereIn("id", id).DeleteE()

	return err
}

func UpdateGroup(id any, fields contracts.Fields) error {
	count := models.Groups().Where("id", "!=", id).
		Where("name", fields["name"]).
		Count()
	if count > 0 {
		return errors.New("分组名称")
	}

	_, err := models.Groups().Where("id", id).
		UpdateE(utils.OnlyFields(fields, "name"))

	return err
}

// HasGroupPermission 判断用户是否存在指定分组的权限
func HasGroupPermission(group *models.Group, user *models.User) bool {
	isUserRoleAdmin := user.Role == models.UserRoleAdmin
	boolIsTrue := group.CreatorId == utils.ToInt(user.Id, 0)
	if isUserRoleAdmin || boolIsTrue {
		return true
	}

	return models.UserGroups().
		Where("group_id", group.Id).
		Where("user_id", user.Id).
		Where("status", models.InviteStatusJoined).
		Count() > 0
}
