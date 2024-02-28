package table

import (
	"bytes"
	"fmt"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"unicode"
	"strings"
)

var FuncMap = template.FuncMap{
	"title": strings.Title,
	"keyexist": func(cols map[string]Column, key string) bool {

		if _, ok := cols[key]; ok {
			return true
		}
		return false

	},
	"tomap": func(from interface{}) map[string]interface{} {

		if structs.IsStruct(from) {
			out := structs.Map(from)
			return out
		} else {
			return from.(map[string]interface{})
		}
	},
	"parse": func(key string, cols []Column, row map[string]interface{}) template.HTML {

		rowtmpl := fmt.Sprintf("{{.%s}}", key)

		if val, ok := GetColumnFromKey(cols, key); ok {

			if len(val.Template) > 0 {

				rowtmpl = val.Template

			}

			if val.Callback != nil {

				return val.Callback(row)

			}

		}

		tmpl, err := template.New("col").Parse(rowtmpl)

		if err != nil {
			return template.HTML(err.Error())
		}

		var w bytes.Buffer
		err2 := tmpl.Execute(&w, row)

		if err2 != nil {

			return template.HTML(err2.Error())
		}

		return template.HTML(w.String())

	},
}

type Table struct {
	Template     *template.Template
	Class        string
	Columns      []Column
	Opts         *Opts
	Builder      *QueryBuilder
	Data         interface{}
	DataRendered string
	Paginator    *Paginator
	Url          string
	Request      *http.Request
}

type Opts struct {
	HideHeader bool
}

func New(opts ...Opts) *Table {

	t := new(Table)

	if len(opts) > 0 {
		t.Opts = &opts[0]
	} else {
		t.Opts = new(Opts)
	}

	return t

}

type Column struct {
	ID       string
	Label    string
	Template string
	Callback func(interface{}) template.HTML
	Ordering string
}

// This is when you have a complex query and you want the same version for row count
func (t *Table) QueryParser(sql string) (string, string) {

	r, _ := regexp.Compile(`(?i)select\s+(.*?)\s*from\s+(.*?)\s*(where\s(.*?)\s*)?`)

	matches := r.FindStringSubmatch(sql)

	var sql2 string

	if len(matches) > 0 {

		sql2 = strings.Replace(sql, matches[1], " COUNT(*) as ct", 1)

	}

	return sql2, sql //original
}

func (t *Table) QueryBuilder(columns []string, table string) *QueryBuilder {

	t.Builder = &QueryBuilder{Columns: columns, Table: table}

	pg := 0
	perpage := 25

	if t.Request !=nil {

	pg, _ = strconv.Atoi(t.Request.URL.Query().Get("page"))
	perpage, _ = strconv.Atoi(t.Request.URL.Query().Get("perpage"))

	}

	if perpage < 1 {
		perpage = 25
	}

	if pg > 1 {
		t.Builder.Start =  (pg - 1) * perpage
	} else {
		t.Builder.Start = 0
	}

	t.Builder.Limit = perpage

	t.Builder.OrderDir = t.Request.URL.Query().Get("order")
	t.Builder.OrderColumn = t.Request.URL.Query().Get("order.column")

	return t.Builder

}

func (t *Table) SetRequest(r *http.Request) {

	t.Request = r
	t.SetUrl(t.Request.URL.String())

}

func (t *Table) SetUrl(u string) {

	t.Url = u

}

func (c Column) Render() template.HTML {

	if len(c.Ordering) > 0 {

		stat := QueryKey(c.Ordering, "order")

		var ud string

		switch stat {
		case "DESC":
			ud = "up"
		case "ASC":
			ud = "down"
		}

		return template.HTML(fmt.Sprintf(`<a class="sorting" href="%s" ><i class="fa fa-angle-%s" ></i>&nbsp;%s</a>`, c.Ordering, ud, c.Label))
	}

	return template.HTML(fmt.Sprintf("%s", c.Label))

}

func (t *Table) SetPaginator(total int, perpage_current ...int) *Paginator {

	var perpage int
	var current int

	if len(perpage_current) > 1 {
		perpage = perpage_current[0]
		current = perpage_current[1]
	}

	if t.Request != nil {
		pg, _ := strconv.Atoi(t.Request.URL.Query().Get("page"))
		ppg, _ := strconv.Atoi(t.Request.URL.Query().Get("perpage"))
		perpage = ppg
		current = pg
	}

	if perpage < 1 {
		perpage = 25
	}

	if current < 1 {
		current = 1
	}

	t.Paginator = &Paginator{Total: total, Perpage: perpage, Current: current}
	t.Paginator.SetUrl(t.Url)

	return t.Paginator

}

func TableTemplate(my_template string) { 

 default_template = my_template

}

func PaginationTemplate(my_template string) { 

 default_pagination = my_template
 
}

func (t *Table) Render() (string, error) {

	final_template := default_template
	var err error
	var tmpl *template.Template

	if t.Template != nil {

		tmpl = t.Template

	} else {

		tmpl, err = template.New("table").Funcs(FuncMap).Parse(final_template)

	}

	if err != nil {

		return "", err
	}

	var w bytes.Buffer
	err2 := tmpl.Execute(&w, t)

	if err2 != nil {

		return "", err2
	}

	return w.String(), nil

}

func (t *Table) AddColumn(name string, id string, templ interface{}, order ...bool) *Table {

	ordering := ""

	if len(order) > 0 {

		key := QueryKey(t.Url, "order")

		if key == "DESC" {

			key = "ASC"
		} else {
			key = "DESC"
		}

		ordering = UpdateUrl(t.Url, "order", key)
		ordering = UpdateUrl(ordering, "order.column",  ToSnakeCase(id))

	}

	tempcol := Column{ID: id, Label: name, Ordering: ordering}

	if reflect.TypeOf(templ).Kind() == reflect.Func {

		tempcol = Column{ID: id, Label: name, Ordering: ordering, Callback: templ.(func(interface{}) template.HTML)}

	}

	var strType = reflect.TypeOf("")

	if reflect.TypeOf(templ) == strType {

		tempcol = Column{ID: id, Label: name, Ordering: ordering, Template: templ.(string)}

	}

	t.Columns = append(t.Columns, tempcol)

	return t
}

func ToSnakeCase(out string) (string) {
	var parts []rune

	spl := strings.Split(out, ".")

	if len(spl) > 1 {
		out = spl[1]
	}

	for k, r := range out {

		if r >= 'A' && r <= 'Z' {
			if k > 0 {
				parts = append(parts, '_')
			}
			rr := unicode.ToLower(r)
			parts = append(parts, rr)
		} else {
			parts = append(parts, r)
		}

	}

	var out2 string

	if len(spl) > 1 {
		out2 = fmt.Sprintf("%s.%s", spl[0], string(parts))
	} else {
		out2 = fmt.Sprint(string(parts))
	}

	return out2
}

func ToInterface(input interface{}, output interface{}) error {

	return mapstructure.Decode(input, output)

}

func UpdateUrl(urlraw string, key string, val string) string {

	urlA, err := url.Parse(urlraw)

	if err != nil {
		return urlraw
	}

	values := urlA.Query()

	values.Set(key, val)

	urlA.RawQuery = values.Encode()

	return urlA.String()

}

func GetColumnFromKey(c []Column, key string) (Column, bool) {

	for _, v := range c {

		if v.ID == key {

			return v, true

		}

	}

	return Column{}, false
}

func QueryKey(urlraw string, key string) string {

	urlA, err := url.Parse(urlraw)

	if err != nil {
		return ""
	}

	values := urlA.Query()

	return values.Get(key)

}
