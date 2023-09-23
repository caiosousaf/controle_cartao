package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
)

type Parametros struct {
	Campos      []string
	OrderCampo  string
	CamposNome  map[string]string
	Limite      uint64
	Offset      uint64
	Filtros     map[string][]string
	OrderByNome bool
	Desc        bool
	Total       bool
	Aggregate   bool
	Chart       bool
}

const (
	sourceTag      = "apelido"
	destinationTag = "sql"
	aliasTag       = "alias"
	distinctTag    = "distinct"
)

func ParseParams(c *gin.Context) (parametros Parametros, err error) {
	lim, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		return
	}

	const MaxLimit = 10000
	if lim <= 0 {
		lim = MaxLimit
	}
	parametros.Limite = uint64(Min(lim, MaxLimit))

	off, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		return
	}
	parametros.Offset = uint64(off)

	parametros.Campos, _ = c.GetQueryArray("campo")

	parametros.OrderCampo = c.DefaultQuery("order", "")

	parametros.Desc, err = strconv.ParseBool(c.DefaultQuery("desc", "false"))
	if err != nil {
		return
	}

	if parametros.Total, err = strconv.ParseBool(c.DefaultQuery("total", "false")); err != nil {
		return
	}

	parametros.Filtros = map[string][]string{}
	for k, v := range c.Request.URL.Query() {
		if k == "limit" || k == "offset" || k == "order" || k == "campo" || k == "desc" {
			continue
		}

		if len(v) > 0 {
			parametros.Filtros[k] = append(parametros.Filtros[k], v...)
		}
	}

	return
}

func (p *Parametros) ValidarOrdenador(dst interface{}) string {
	const ()

	elemDst := reflect.ValueOf(dst).Elem()

	if elemDst.Kind() != reflect.Struct || elemDst.NumField() == 0 {
		return ""
	}

	var orderBy string

	for s := 0; s < elemDst.NumField(); s++ {
		field := elemDst.Type().Field(s)
		inTag := field.Tag.Get(sourceTag)
		outTag := field.Tag.Get(destinationTag)
		alias := field.Tag.Get(aliasTag)

		if inTag == "" || outTag == "" {
			continue
		}

		if p.OrderCampo == inTag {
			if p.OrderByNome {
				orderBy = inTag
			} else {
				if alias != "" {
					orderBy = alias + "." + outTag
				} else {
					orderBy = outTag
				}
			}
			break
		}

		// Se nenhum campo corresponder ao campo de ordenação especificado,
		// continue procurando para usar o primeiro campo encontrado como padrão.
		if orderBy == "" {
			if p.OrderByNome {
				orderBy = inTag
			} else {
				if alias != "" {
					orderBy = alias + "." + outTag
				} else {
					orderBy = outTag
				}
			}
		}
	}

	if p.Desc {
		orderBy += " DESC"
	} else {
		orderBy += " ASC"
	}

	return orderBy
}

func (p *Parametros) ValidFields(dst interface{}, options ...map[string]string) (selectedFields []string, fieldValues []interface{}, err error) {

	var (
		enableDistinct bool
	)

	p.CamposNome = make(map[string]string)
	elemDst := reflect.ValueOf(dst).Elem()

	if elemDst.Kind() != reflect.Struct {
		err = errors.New("O destino não é uma estrutura válida")
		return
	}

	if elemDst.NumField() == 0 {
		err = errors.New("Nenhum campo disponível na estrutura")
		return
	}

	if p.Total {
		count := "count(1)"
		if len(options) > 0 && options[0]["count"] != "" {
			count = options[0]["count"]
		}
		selectedFields = append(selectedFields, count)
		return
	}

	if len(options) > 0 && options[0]["distinct"] != "" {
		enableDistinct = true
	}

	applyField := func(tagToMatch string) {
		for s := 0; s < elemDst.NumField(); s++ {
			field := elemDst.Type().Field(s)
			sourceFieldTag := field.Tag.Get(sourceTag)
			destinationFieldTag := field.Tag.Get(destinationTag)
			aliasFieldTag := field.Tag.Get(aliasTag)
			distinctFieldTag := field.Tag.Get(distinctTag)
			if sourceFieldTag == "" || destinationFieldTag == "" {
				continue
			}
			internal := ""

			if tagToMatch == sourceFieldTag || tagToMatch == "" {
				pointerToField := reflect.New(reflect.PtrTo(elemDst.Field(s).Type()))
				pointerToField.Elem().Set(elemDst.Field(s).Addr())

				if aliasFieldTag != "" {
					internal = aliasFieldTag + "." + destinationFieldTag
					destinationFieldTag = internal + " AS " + sourceFieldTag
				} else {
					destinationFieldTag = destinationFieldTag + " AS " + sourceFieldTag
				}

				if distinctFieldTag != "" && enableDistinct {
					destinationFieldTag = "DISTINCT ON (" + internal + ") " + destinationFieldTag
				}

				selectedFields = append(selectedFields, destinationFieldTag)
				p.CamposNome[destinationFieldTag] = sourceFieldTag
				fieldValues = append(fieldValues, pointerToField.Elem().Interface())

			}
		}
	}

	if len(p.Campos) > 0 {
		for _, tag := range p.Campos {
			applyField(tag)
		}
	} else {
		applyField("")
	}

	return
}
