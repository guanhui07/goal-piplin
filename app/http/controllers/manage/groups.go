package manage

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
	"github.com/goal-web/querybuilder"
	"github.com/qbhy/goal-piplin/app/models"
	"github.com/qbhy/goal-piplin/app/usecase"
)

func GetGroups(request contracts.HttpRequest, guard contracts.Guard) any {
	user := guard.User().(*models.User)
	name := request.GetString("name")
	data := models.Groups().
		When(user.Role != "admin", func(q contracts.Query[models.Group]) contracts.Query[models.Group] {
			return q.WhereFunc(func(q contracts.Query[models.Group]) {
				q.Where("creator_id", user.Id).OrWhereExists(func() contracts.Query[models.Group] {
					return querybuilder.New[models.Group]("user_groups").
						Where("user_id", user.Id).
						Where("status", models.InviteStatusJoined).
						Where("`groups`.id", querybuilder.Expression("user_groups.group_id"))
				})
			})
		}).
		When(name != "", func(q contracts.Query[models.Group]) contracts.Query[models.Group] {
			return q.Where("name", "like", "%"+name+"%")
		}).
		Get().ToArray()
	return contracts.Fields{
		"data": data,
	}
}

func GetGroupMembers(request contracts.HttpRequest) any {
	group := models.Groups().FindOrFail(request.Get("id"))

	data := table.ArrayQuery("user_groups").
		OrderByDesc("user_groups.id").
		Select("user_id", "username", "nickname", "avatar", "user_groups.created_at", "user_groups.id", "status").
		LeftJoin("`users`", "`users`.id", "=", "user_groups.user_id").
		Where("group_id", group.Id).
		Get().
		ToArray()
	return contracts.Fields{
		"data": data,
	}
}

func CreateGroup(request contracts.HttpRequest, guard contracts.Guard) any {
	data := models.Groups().Create(contracts.Fields{
		"name":       request.GetString("name"),
		"creator_id": guard.User().GetId(),
	})
	return contracts.Fields{
		"data": data,
	}
}

func DeleteGroups(request contracts.HttpRequest) any {
	err := usecase.DeleteGroups(request.Get("id"))

	if err != nil {
		return contracts.Fields{"msg": err.Error()}
	}

	return contracts.Fields{"data": nil}
}

func UpdateGroup(request contracts.HttpRequest) any {
	id := request.Get("id")
	err := usecase.UpdateGroup(id, request.Fields())

	if err != nil {
		return contracts.Fields{"msg": err.Error()}
	}

	return contracts.Fields{"data": nil}
}
