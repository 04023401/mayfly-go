package vo

import "time"

type DbListVO struct {
	Id       *int64  `json:"id"`
	Name     *string `json:"name"`
	Database *string `json:"database"`
	Remark   *string `json:"remark"`
	TagId    *int64  `json:"tagId"`
	TagPath  *string `json:"tagPath"`

	InstanceId   *int64  `json:"instanceId"`
	InstanceName *string `json:"instanceName"`
	InstanceType *string `json:"type"`
	Host         string  `json:"host"`
	Port         int     `json:"port"`
	Username     string  `json:"username"`

	CreateTime *time.Time `json:"createTime"`
	Creator    *string    `json:"creator"`
	CreatorId  *int64     `json:"creatorId"`
	UpdateTime *time.Time `json:"updateTime"`
	Modifier   *string    `json:"modifier"`
	ModifierId *int64     `json:"modifierId"`
}
