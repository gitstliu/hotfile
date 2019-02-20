package vo

type UserInfo struct {
	UserId   int64
	Username string
	Groups   map[string]*UserGroup
}

type UserGroup struct {
	GroupId     int64
	GroupName   string
	TotalSize   int64
	CurrentSize int64
}
