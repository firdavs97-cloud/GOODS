package utils

import "fmt"

const (
	GetRecord = "get_record#"
)

func GetRecordCacheKey(id, projectId int) string {
	return GetRecord + fmt.Sprintf("%d#%d", id, projectId)
}
