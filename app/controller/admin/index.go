package admin

import (
	"gfast/app/service/admin/auth_service"
	"gfast/app/service/admin/user_service"
	"gfast/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type Index struct{}

//后台首页接口数据
func (c *Index) Index(r *ghttp.Request) {
	//获取用户信息
	userInfo := user_service.GetLoginAdminInfo(r)
	//菜单列表
	var menuList g.List
	if userInfo != nil {
		userId := gconv.Int(userInfo["id"])
		delete(userInfo, "user_password")
		//获取用户角色信息
		allRoles, err := auth_service.GetRoleList()
		if err == nil {
			roles, err := user_service.GetAdminRole(userId, allRoles)
			if err == nil {
				name := make([]string, len(roles))
				roleIds := make([]int, len(roles))
				for k, v := range roles {
					name[k] = v.Name
					roleIds[k] = v.Id
				}
				userInfo["roles"] = strings.Join(name, "，")
				//获取菜单信息
				menuList, err = user_service.GetAdminMenusByRoleIds(roleIds)
				if err != nil {
					g.Log().Error(err)
				}
			} else {
				g.Log().Error(err)
				userInfo["roles"] = ""
			}
		} else {
			g.Log().Error(err)
			userInfo["roles"] = ""
		}

	}

	result := g.Map{
		"userInfo": userInfo,
		"menuList": menuList,
	}
	response.SusJson(true, r, "ok", result)
}