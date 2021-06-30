package mysql

import (
	"strings"
	"time"

	amazonHelper "github.com/Jleagle/canihave/pkg/amazon"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/ngs/go-amazon-product-advertising-api/amazon"
	"go.uber.org/zap"
)

type Relation struct {
	ID          string
	RelatedID   string
	DateCreated string
	Type        string
}

func saveSimilarItems(id string, region amazon.Region, itemType string) {

	similar, err := amazonHelper.GetSimilarItems(id, region)

	if err != nil && strings.Contains(err.Error(), "AWS.ECommerceService.NoSimilarities") {
		return
	} else if err != nil {
		logger.Logger.Fatal("", zap.Error(err))
	}

	for _, amazonItem := range similar.Items.Item {

		item := amazonItemToItem(amazonItem, itemType, region)

		saveToMysql(item)
		saveToMemcache(item)

		// Save the relation
		builder := squirrel.Insert("relations")
		builder = builder.Columns("id", "relatedId", "dateCreated", "type")
		builder = builder.Values(id, item.ID, time.Now().Unix(), TypeSimilar)

		err := mysql.Insert(builder)

		if sqlerr, ok := err.(*mysql.MySQLError); ok {
			if sqlerr.Number == 1062 { // Duplicate entry
				continue
			}
		}

		if err != nil {
			logger.Logger.Error("Can't insert related item", zap.Error(err))
		}
	}
}

func (i Item) GetRelated(itemType string) (items []Item) {

	builder := squirrel.Select("relatedId").From("relations").Where("id = ? AND type = ?", i.ID, itemType)
	rows := mysql.Query(builder)
	defer rows.Close()

	var ids []string
	for rows.Next() {

		var id string
		err := rows.Scan(&id)
		if err != nil {
			logger2.Err("Can't scan related item", zap.Error(err))
		}
		ids = append(ids, id)
	}

	return GetMulti(ids, i.Region, itemType)
}
