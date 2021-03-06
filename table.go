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
	"strings"
)



var funcMap = template.FuncMap {
 "title": strings.Title,
 "keyexist": func(cols map[string]Column,key string) (bool){

 	if _, ok := cols[key]; ok {
 		return true
	}
	return false

 },
 "tomap": func(from interface{}) (map[string]interface{}){

    out := structs.Map(from)
    return out
 },
  "parse": func(key string, cols []Column, row map[string]interface{}) (template.HTML) {

       rowtmpl := fmt.Sprintf("{{.%s}}",key)

       if val,ok := GetColumnFromKey(cols,key); ok {

       	  if(len(val.Template)>0){

             rowtmpl = val.Template

       	  }

       	  if(val.Callback != nil){

       	  	return val.Callback(row)

		  }

	   }

	  tmpl, err := template.New("col").Parse(rowtmpl)

	  if(err!=nil){
		  return template.HTML(err.Error())
	  }

	  var w bytes.Buffer
	  err2 := tmpl.Execute(&w,row)

	  if(err2!=nil){

		  return template.HTML(err2.Error())
	  }

	  return template.HTML(w.String())

  },
 }

type Table struct {
	Template *template.Template
	Class string
	Columns []Column
	Data interface{}
	DataRendered string
	Paginator *Paginator
	Url string
	Request *http.Request
}

func New() *Table {

	t := new(Table)

	return t

}

type Column struct {
	ID string
	Label string
	Template string
	Callback func (interface{}) template.HTML
	Ordering string
}


func (t *Table) SetRequest(r *http.Request){

	t.Request = r
	t.SetUrl(t.Request.URL.String())

}

func(t *Table) SetUrl(u string) {

	t.Url = u

}

func(c Column) Render() (template.HTML) {

	if len(c.Ordering)>0 {

		stat := QueryKey(c.Ordering,"order")

		var ud string

		switch stat {
		 case "DESC":
		    ud = "up"
		 case "ASC":
		 	ud = "down"
		}

		return template.HTML(fmt.Sprintf(`<a class="sorting" href="%s" ><i class="fa fa-angle-%s" ></i>&nbsp;%s</a>`,c.Ordering,ud, c.Label))
	}

	return template.HTML(fmt.Sprintf("%s",c.Label))

}



func (t *Table) SetPaginator(total int, perpage int,current int) (*Paginator){

    t.Paginator = &Paginator{Total: total,Perpage: perpage,Current: current}
    t.Paginator.SetUrl(t.Url)

    return t.Paginator

}

func (t *Table) Render() (string,error){

	final_template := default_template
    var err error
	var tmpl *template.Template

	if(t.Template!=nil) {

		tmpl = t.Template
		tmpl = tmpl.Funcs(funcMap) //add in built in funcs

	}else{

	    tmpl, err = template.New("table").Funcs(funcMap).Parse(final_template)

	}

	if(err!=nil){

		return "",err
	}

    var w bytes.Buffer
	err2 := tmpl.Execute(&w,t)

	if(err2!=nil){

		return "",err2
	}

	return w.String(),nil

}

func (t *Table) AddColumn(name string,id string,templ interface{},order ...bool) (*Table){

	ordering := ""

	if( len(order)>0 ){

		key := QueryKey(t.Url,"order")

		if(key=="DESC"){

			key = "ASC"
		}else {
			key = "DESC"
		}

		ordering = UpdateUrl(t.Url,"order",key)
		ordering = UpdateUrl(ordering,"order.column",ToSnakeCase(id))

	}

	tempcol := Column{ID: id, Label: name, Ordering: ordering}

	if reflect.TypeOf(templ).Kind() == reflect.Func {

		tempcol = Column{ID: id,Label: name,  Ordering: ordering, Callback: templ.(func (interface{}) (template.HTML)) }

	}

	var strType = reflect.TypeOf("")

	if reflect.TypeOf(templ) == strType {

		tempcol = Column{ID: id, Label: name,  Ordering: ordering, Template: templ.(string)}

	}

	t.Columns = append(t.Columns,tempcol)

	return t
}



func ToSnakeCase(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// irregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && ( // head and tail
			(i > 0 && rune(camel[i-1]) >= 'a') || // pre
				(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

func ToInterface(input interface{},output interface{}) (error){

	return mapstructure.Decode(input,output)

}

func UpdateUrl(urlraw string,key string, val string) (string){

	urlA, err := url.Parse(urlraw)

	if err != nil {
		return urlraw
	}

	values := urlA.Query()

	values.Set(key, val)

	urlA.RawQuery = values.Encode()

	return urlA.String()

}

func GetColumnFromKey(c []Column,key string)(Column,bool){

	for _,v := range c {

		if(v.ID==key){

			return v,true

		}

	}

	return Column{},false
}

func QueryKey(urlraw string,key string) (string){

	urlA, err := url.Parse(urlraw)

	if err != nil {
		return ""
	}

	values := urlA.Query()

	return values.Get(key)

}
