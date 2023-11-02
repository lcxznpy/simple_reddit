package mysql

import (
	"database/sql"
	"goweb/models"

	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	sqlStr := `select community_id,community_name from community `
	if err := db.Select(&data, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return

}

// GetCommunityDetailById 根据id查询社区详情
func GetCommunityDetailById(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)

	sqlStr := `select 
    		community_id,community_name,introduction,create_time 
			from community
			where community_id = ?`

	err = db.Get(community, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidId
		}
	}
	return
}
