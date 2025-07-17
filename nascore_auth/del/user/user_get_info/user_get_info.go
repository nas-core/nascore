package user_get_info

import "strconv"

type User struct {
	Username string
	HomeDir  string
	IsAdmin  bool
	UserId   int64
}

// getUserInfo 获取用户信息和权限
func GetUserInfo(username string, allusers []map[string]string) (user User, err error) {
	for _, u := range allusers {
		if u["username"] == username {
			user.Username = u["username"]
			user.HomeDir = u["home"]
			if u["isadmin"] == "yes" || u["isadmin"] == "true" || u["isadmin"] == "1" {
				user.IsAdmin = true
			} else {
				user.IsAdmin = false
			}
			id, err := strconv.ParseInt(u["id"], 10, 64)
			if err != nil {
				user.UserId = 0
			}
			user.UserId = id
			return user, nil
		}
	}
	return user, nil
}
