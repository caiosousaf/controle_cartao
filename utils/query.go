package utils

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"reflect"
)

func ConfigurarPaginacao(p *Parametros, model interface{}, query *sq.SelectBuilder, possuiOrdenador ...bool) (result interface{}, next *bool, count *int64, err error) {
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	slice := reflect.MakeSlice(reflect.SliceOf(modelType), 0, 0)
	if p.Total {
		var total int64
		var rows *sql.Rows
		rows, err = query.
			Query()

		if err != nil {
			return
		}

		for rows.Next() {
			total++
		}

		count = &total
	} else {

		pre := query.
			Limit(p.Limite + 1).
			Offset(p.Offset)

		// é necessário verificar se o erdenador é pre-selecionado e
		// adiciona-lo caso não seja
		if len(possuiOrdenador) != 1 || !possuiOrdenador[0] {
			pre = pre.OrderBy(p.ValidarOrdenador(model))
		}

		rows, err := pre.Query()
		if err != nil {
			return result, next, count, err
		}

		_, values, err := p.ValidFields(model)
		if err != nil {
			return result, next, count, err
		}

		slice = reflect.MakeSlice(reflect.SliceOf(modelType), 0, int(p.Limite+1))
		for rows.Next() {
			if err = rows.Scan(values...); err != nil {
				return result, next, count, err
			}
			slice = reflect.Append(slice, reflect.Indirect(reflect.ValueOf(model)))
		}

		hasNext := slice.Len() > int(p.Limite)
		if hasNext {
			slice = slice.Slice(0, int(p.Limite))
		}

		next = &hasNext
	}

	result = slice.Interface()

	return
}

// SelectBuilderToString recebe um squirrel.SelectBuilder e retorna ela em string
func SelectBuilderToString(builder sq.SelectBuilder) (build string, err error) {
	build, _, err = builder.ToSql()
	if err != nil {
		return build, err
	}

	return
}
