package manage

import (
	"github.com/goal-web/contracts"
	"github.com/qbhy/goal-piplin/app/models"
	"github.com/qbhy/goal-piplin/app/usecase"
)

func GetCabinets(request contracts.HttpRequest, guard contracts.Guard) any {
	user := guard.User().(*models.User)
	perPageSize := request.Int64Optional("perPageSize", 10)
	page := request.Int64Optional("current", 1)
	name := request.GetString("name")

	list, total := models.Cabinets().
		OrderByDesc("id").
		When(name != "", func(q contracts.Query[models.Cabinet]) contracts.Query[models.Cabinet] {
			return q.Where("name", "like", "%"+name+"%")
		}).
		When(user.Role != "admin", func(q contracts.Query[models.Cabinet]) contracts.Query[models.Cabinet] {
			return q.Where("creator_id", user.Id)
		}).
		Paginate(perPageSize, page)
	return contracts.Fields{
		"total": total,
		"data":  list.ToArray(),
	}
}

func CreateCabinet(request contracts.HttpRequest, guard contracts.Guard) any {
	name := request.GetString("name")
	settings := request.Get("settings")
	cabinet, err := usecase.CreateCabinet(guard.GetId(), name, settings)

	if err != nil {
		return contracts.Fields{"msg": err.Error()}
	}

	return contracts.Fields{"data": cabinet}
}

func UpdateCabinet(request contracts.HttpRequest) any {
	id := request.Get("id")
	name := request.GetString("name")
	settings := request.Get("settings")
	err := usecase.UpdateCabinet(id, name, settings)

	if err != nil {
		return contracts.Fields{"msg": err.Error()}
	}

	return contracts.Fields{"data": nil}
}

func DeleteCabinet(request contracts.HttpRequest) any {
	id := request.Get("id")
	err := usecase.DeleteCabinet(id)

	if err != nil {
		return contracts.Fields{"msg": err.Error()}
	}

	return contracts.Fields{"data": nil}
}
