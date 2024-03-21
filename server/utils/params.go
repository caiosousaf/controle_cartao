package utils

import (
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reflect"
	"strconv"
)

// FlagFiltro é usado para definir o tipo do filtro
type FlagFiltro int

// Parametros é usado para requisições quando o parametro Query é necessário
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

// Filtro é a representação de um filtro disponível
type Filtro struct {
	Valor   string
	Flag    FlagFiltro
	Tamanho int
}

const (
	sourceTag      = "apelido"
	destinationTag = "sql"
	// aliasTag é uma tag que serve para dizer qual é alias do campo
	aliasTag = "alias"
	// distinctTag é uma tag especifica para se colocar distinct em um campo
	distinctTag = "distinct"
	// MaxLimit define um valor máximo que uma listagem pode requisitar
	MaxLimit = 10000
)

const (
	FlagFiltroNenhum FlagFiltro = 1 << iota
	FlagFiltroEq
	FlagFiltroIn
	FlagFiltroNotIn
)

// CriarFiltros cria os filtros
func CriarFiltros(v string, flag FlagFiltro, tamanho ...int) Filtro {
	if len(tamanho) == 0 {
		tamanho = []int{1}
	}

	return Filtro{
		Valor:   v,
		Flag:    flag,
		Tamanho: tamanho[0],
	}
}

// CriarFiltros retorna um squirrel.SelectBuilder com todos os filtros aplicados a ele
func (p *Parametros) CriarFiltros(builder sq.SelectBuilder, disponiveis map[string]Filtro) sq.SelectBuilder {
	for k := range disponiveis {
		var v = disponiveis[k]
		for k1, v1 := range p.Filtros {
			if k == k1 {
				v.Valor = "( " + v.Valor + " )"
				switch v.Flag {
				case FlagFiltroIn:
					builder = builder.Where(sq.Eq{
						v.Valor: v1,
					})
				case FlagFiltroNotIn:
					builder = builder.Where(sq.NotEq{
						v.Valor: v1,
					})
				case FlagFiltroEq:
					builder = builder.Where(v.Valor, func(xs []string) (v []interface{}) {
						for x := range xs {
							v = append(v, xs[x])
						}
						return
					}(v1[0:v.Tamanho])...)
				}
			}
		}
	}

	return builder
}

func ParseParams(c *gin.Context) (parametros Parametros, err error) {
	lim, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		return
	}

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

// GetUUIDFromParam pega um parametro da rota e retorna um ponteiro de uuid
func GetUUIDFromParam(c *gin.Context, paramName string) (*uuid.UUID, error) {
	id := c.Param(paramName)
	uuidObj, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return &uuidObj, nil
}

// TemFiltro verifica se o filtro enviado existe
func (p *Parametros) TemFiltro(f string) (temFiltro bool) {
	if _, existe := p.Filtros[f]; existe {
		return true
	}

	return
}

// AdicionarFiltro adiciona um filtro para Parametros
func (p *Parametros) AdicionarFiltro(nomeFiltro string, values ...string) *Parametros {
	if len(values) > 0 {
		if p.Filtros == nil {
			p.Filtros = map[string][]string{}
		}
		p.Filtros[nomeFiltro] = values
	}

	return p
}

// LimparFiltros limpa todos os filtros dos Parametros
func (p *Parametros) LimparFiltros() *Parametros {
	p.Filtros = make(map[string][]string)

	return p
}

// RemoverFiltros limpara um ou todos os filtros dos parametros
func (p *Parametros) RemoverFiltros(f ...string) *Parametros {
	for i := range f {
		delete(p.Filtros, f[i])
	}

	return p
}
