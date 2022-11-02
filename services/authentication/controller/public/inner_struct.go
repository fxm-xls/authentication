package public

type DeptAdd struct {
	DeptName  string `json:"dept_name"`
	ChargerId int    `json:"charger_id"`
	RoleId    []int  `json:"role_id"`
	ParentId  int    `json:"parent_id"`
	DeptLevel int    `json:"dept_level"`
}

type DeptUpdate struct {
	DeptId       int    `json:"dept_id"`
	DeptName     string `json:"dept_name"`
	DeptNameNew  string `json:"dept_name_new"`
	ChargerId    int    `json:"charger_id"`
	ChargerIdNew int    `json:"charger_id_new"`
	RoleId       []int  `json:"role_id"`
	RoleIdNew    []int  `json:"role_id_new"`
}

type UserAdd struct {
	Account      string `json:"account"`
	UserName     string `json:"user_name"`
	DeptId       int    `json:"dept_id"`
	Pwd          string `json:"pwd"`
	RoleId       int    `json:"role_id"`
	ManageDeptId []int  `json:"manage_dept_id"`
}

type UserUpdate struct {
	UserId          int    `json:"user_id"`
	UserName        string `json:"user_name"`
	UserNameNew     string `json:"user_name_new"`
	RoleId          int    `json:"role_id"`
	RoleIdNew       int    `json:"role_id_new"`
	ManageDeptId    []int  `json:"manage_dept_id"`
	ManageDeptIdNew []int  `json:"manage_dept_id_new"`
}
