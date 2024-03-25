package good

import (
	"encoding/json"
	"errors"
	"fmt"
	"goods/models/utils"
	"goods/pkg/cache"
	"goods/pkg/postgres/db"
	"strings"
	"time"
)

const (
	tbName = "GOODS"
)

type Good struct {
	Id          int       `json:"id" db:"id" pk:"true"`
	ProjectId   int       `json:"projectId" db:"project_id" fk:"true"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Priority    int       `json:"priority" db:"priority" au:"true"`
	Removed     bool      `json:"removed" db:"removed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" au:"true"`
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (g *Good) MarshalBinary() ([]byte, error) {
	// Marshal the object to JSON
	return json.Marshal(g)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (g *Good) UnmarshalBinary(data []byte) error {
	// Unmarshal JSON data into the object
	return json.Unmarshal(data, g)
}

func New(name string, projectId int) (Good, error) {
	good := Good{
		ProjectId:   projectId,
		Name:        name,
		Description: "Описание",
	}
	c, v := utils.GetColumnsAndValues(good, "pk", "au")
	err := db.InsertRecord(tbName, c, v, []string{"id", "priority", "created_at"}, &good.Id, &good.Priority, &good.CreatedAt)
	if err == nil {
		cache.Set(utils.GetRecordCacheKey(good.Id, good.ProjectId), good, 0)
	}
	return good, err
}

func Get(id, projectId int) (*Good, error) {
	var g Good
	rowsStr, err := cache.Get(utils.GetRecordCacheKey(id, projectId))
	if err == nil {
		err = json.Unmarshal([]byte(rowsStr), &g)
	}
	if err == nil {
		return &g, nil
	}
	rows, err := db.GetRecord(tbName, db.QueryConstruct{
		WhereExpr:       "id=$1 AND project_id=$2",
		WhereExprValues: []interface{}{id, projectId},
	})
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errors.New("item not found")
	}
	err = rows.Scan(&g.Id, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

type Report struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
}

func ListRecordsAndReport(limit, offset int) (*Report, []Good, error) {
	rows, err := db.GetRecord(tbName, db.QueryConstruct{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, nil, err
	}
	r := Report{}
	goods := make([]Good, 0)
	for rows.Next() {
		r.Total++
		g := Good{}
		err = rows.Scan(&g.Id, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
		if err != nil {
			return nil, nil, err
		}
		if g.Removed {
			r.Removed++
		}
		goods = append(goods, g)
	}
	return &r, goods, nil
}

func Remove(id, projectId int, old *Good) error {
	var err error
	if old == nil {
		old, err = Get(id, projectId)
		if err != nil {
			return err
		}
	}
	q := db.QueryConstruct{
		WhereExpr:       "id=$1 AND project_id=$2",
		WhereExprValues: []interface{}{id, projectId},
		SetExpr:         "removed=true",
	}
	err = db.UpdateRecord(tbName, q)
	if err == nil {
		cache.Delete(utils.GetRecordCacheKey(id, projectId))
	}
	return err
}

func (g *Good) Save(old *Good) error {
	var err error
	if old == nil {
		old, err = Get(g.Id, g.ProjectId)
		if err != nil {
			return err
		}
	}
	q := db.QueryConstruct{
		WhereExpr:       "id=$1 AND project_id=$2",
		WhereExprValues: []interface{}{g.Id, g.ProjectId},
		SetExprValues:   make([]interface{}, 0, 5),
	}
	k := 1
	q.SetExpr += fmt.Sprintf("name=$%d,", k)
	q.SetExprValues = append(q.SetExprValues, g.Name)
	if g.Description != old.Description {
		k++
		q.SetExpr += fmt.Sprintf("description=$%d,", k)
		q.SetExprValues = append(q.SetExprValues, g.Description)
	}
	q.WhereExpr = fmt.Sprintf("id=$%d AND project_id=$%d", k+1, k+2)
	if strings.HasSuffix(q.SetExpr, ",") {
		q.SetExpr = q.SetExpr[:len(q.SetExpr)-1]
	}
	err = db.UpdateRecord(tbName, q)
	if err == nil {
		cache.Set(utils.GetRecordCacheKey(g.Id, g.ProjectId), g, 0)
	}
	return err
}

type UpdatedPriorities struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
}

func (g *Good) Reprioritiize(newPriority int) ([]UpdatedPriorities, error) {
	q1 := db.QueryConstruct{
		WhereExpr:       "id=$2 AND project_id=$3",
		WhereExprValues: []interface{}{g.Id, g.ProjectId},
		SetExpr:         "priority=$1",
		SetExprValues:   []interface{}{g.Priority},
	}
	q2 := db.QueryConstruct{
		WhereExpr:       "priority>=$1 AND priority<$2",
		WhereExprValues: []interface{}{newPriority, g.Priority},
		SetExpr:         "priority = priority + 1",
	}
	if newPriority > g.Priority {
		q2 = db.QueryConstruct{
			WhereExpr:       "priority>$1 AND priority<=$2",
			WhereExprValues: []interface{}{g.Priority, newPriority},
			SetExpr:         "priority = priority - 1",
		}
	}

	g.Priority = newPriority
	err := db.TransactUpdate(tbName, q2, q1)
	if err != nil {
		return nil, err
	}
	cache.FlushAll()

	q2.WhereExpr = "priority>=$1 AND priority<=$2"
	q2.Columns = "id, priority"

	rows, err := db.GetRecord(tbName, q2)
	if err != nil {
		return nil, err
	}

	c := newPriority - g.Priority
	if c < 0 {
		c = -c
	}

	res := make([]UpdatedPriorities, c, c)

	for i, _ := range res {
		if rows.Next() {
			err = rows.Scan(&res[i].ID, &res[i].Priority)
		} else {
			err = errors.New("db error on changing priorities")
		}
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
