package table

import (
	"fmt"
	"strings"
)

var QueryMode string = "mysql"

type QueryBuilder struct {
	Columns []string
	Table   string
	Params  []string
	Args    []interface{}
	ArgsRaw string
	Start   int
	Limit   int
	OrderColumn string
	OrderDir string
}

func (b *QueryBuilder) Arg(column, operator, value interface{}) *QueryBuilder {

	var item string

	if QueryMode == "postgres" {
		index := len(b.Params)
		item = fmt.Sprintf("%s %s $%d", column, operator, index)
	} else {
		item = fmt.Sprintf("%s %s ?", column, operator)
	}

	b.Params = append(b.Params, item)
	b.Args = append(b.Args, value)

	return b

}

func (b *QueryBuilder) GetArgs() []interface{} {
	return b.Args
}

/*

var sql_part []string
	sql_begin := fmt.Sprintf("SELECT %s FROM %s WHERE", columns, table)
	sql_part = append(sql_part, sql_begin)

*/

func (b *QueryBuilder) Fetch() (string, string) {

	var sql1_part []string
	var sql2_part []string

	sql1 := fmt.Sprintf("SELECT COUNT(*) FROM %s ", b.Table)

	sql1_part = append(sql1_part, sql1)

	if len(b.ArgsRaw) < 1 {
		b.ArgsRaw = strings.Join(b.Params, " AND ")
	}

	if len(b.ArgsRaw) > 0 {
		sql1_part = append(sql1_part, fmt.Sprintf("WHERE %s", b.ArgsRaw))
	}

	sql2 := fmt.Sprintf("SELECT %s FROM %s ", strings.Join(b.Columns, ","), b.Table)

	sql2_part = append(sql2_part, sql2)

	if len(b.ArgsRaw) > 0 {
		sql2_part = append(sql2_part, fmt.Sprintf("WHERE %s", b.ArgsRaw))
	}

	if len(b.OrderDir)>0 && len(b.OrderColumn)>0 {

		sql2_part = append(sql2_part, fmt.Sprintf("ORDER BY %s %s", b.OrderColumn,b.OrderDir))

	}

	if b.Start >= 0 && b.Limit > 0 {

		if QueryMode == "postgres" {

			sql2_part = append(sql2_part, fmt.Sprintf("OFFSET %d LIMIT %d", b.Start, b.Limit))

		}else{

			sql2_part = append(sql2_part, fmt.Sprintf("LIMIT %d,%d", b.Start, b.Limit))

		}
	}

	return strings.Join(sql1_part, " "), strings.Join(sql2_part, " ")

}
