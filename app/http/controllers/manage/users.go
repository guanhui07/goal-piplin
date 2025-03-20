package manage

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"github.com/goal-web/validation"
	"github.com/qbhy/goal-piplin/app/http/requests"
	"github.com/qbhy/goal-piplin/app/models"
	"github.com/qbhy/goal-piplin/app/usecase"
)

func GetUsers(request contracts.HttpRequest) any {
	name := request.GetString("name")
	perPage := request.Int64Optional("pageSize", 10)
	page := request.Int64Optional("current", 1)
	list, total := models.Users().
		OrderByDesc("id").
		When(name != "", func(q contracts.Query[models.User]) contracts.Query[models.User] {
			return q.Where("name", "like", "%"+name+"%")
		}).
		Paginate(perPage, page)

	return contracts.Fields{
		"total": total,
		"data":  list.ToArray(),
	}
}

func CreateUser(request requests.CreateUserRequest) any {
	validation.VerifyForm(request)

	userName := request.GetString("username")
	passwd := request.GetString("password")
	role := request.GetString("role")
	newUser, err := usecase.CreateUser(userName, passwd, role)
	if err != nil {
		return contracts.Fields{"msg": "创建用户失败：" + err.Error()}
	}
	return contracts.Fields{
		"data": newUser,
	}
}

func DeleteUsers(request contracts.HttpRequest) any {
	id := request.Get("id")
	err := usecase.DeleteUsers(id)

	if err != nil {
		return contracts.Fields{"msg": err.Error()}
	}

	return contracts.Fields{"data": nil}
}

func UpdateUser(request contracts.HttpRequest, hash contracts.Hasher) any {
	fields := request.Fields()
	if fields["password"] != nil {
		hashPwd := hash.Make(utils.ToString(fields["password"], ""), nil)
		fields["password"] = hashPwd
	}
	id := request.Get("id")
	err := usecase.UpdateUser(id, fields)

	if err != nil {
		return contracts.Fields{"msg": err.Error()}
	}

	return contracts.Fields{"data": nil}
}
