package table

import (
	"bytes"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Paginator struct {
	Total int
	Perpage int
	Current int
	HasPrev bool
	HasNext bool
	Pages []int
	PageLink string
	PageLinkFirst string
	PageLinkPrev string
	PageLinkNext string
	PageLinkLast string
	Url string
	PageList []int
	Template *template.Template
	Request *http.Request
}
func(p *Paginator) SetRequest(r *http.Request) {

	 p.Request = r
	 p.SetUrl(p.Request.URL.String())

	//override the perpage var
	pp := QueryKey(p.Url,"perpage")

	if(len(pp)>0){
		p.Perpage,_ = strconv.Atoi(pp)
	}

	if p.Perpage < 10 {
		p.Perpage = 10
	}

	if p.Perpage >250 {
		p.Perpage = 250
	}

	if p.Current> 1{

		p.Current = 1

	}

}

func(p *Paginator) SetUrl(u string) {

	p.Url = u

}

func (p *Paginator) GetUrl(page int) string {

	pp := strconv.Itoa(p.Perpage)
	return UpdateUrl(UpdateUrl(p.Url, "page", strconv.Itoa(page)), "perpage", pp)

}

func (p *Paginator) GetPerpage(page int) string{

	return UpdateUrl(p.Url,"perpage",strconv.Itoa(page))
}


func (p *Paginator) GetOffsets() (int,int){

	if(p.Perpage == 0){

		p.Perpage = 10
	}

	a  := (p.Current-1) * p.Perpage

	b  := p.Perpage

	return a,b
}


func(p *Paginator) Render() (template.HTML) {


	if(len(p.PageList)<1){
		p.PageList=[]int{10,25,50,100}
	}

	pagination_tmpl := default_pagination

	//make pages
	p.Pages = make([]int,0)

	pgsint := float64(p.Total) /  float64(p.Perpage)

	pgs := int(math.Ceil(pgsint))

	for i:=1; i<=pgs; i++ {

		p.Pages = append(p.Pages,i)

	}

	p.HasPrev = true
	p.HasNext = true

	p.PageLinkPrev = p.GetUrl(p.Current-1)
	p.PageLinkNext = p.GetUrl(p.Current+1)
	p.PageLinkLast = p.GetUrl(pgs)

	if (p.Current <2){
		p.HasPrev = false
		p.PageLinkPrev = p.GetUrl(1)
	}

	if(p.Current == pgs){
		p.HasNext = false
		p.PageLinkNext = p.GetUrl(pgs)
	}

	var tmpl *template.Template
	var err error

	if(p.Template!=nil){

		tmpl = p.Template

	}else{

		tmpl, err = template.New("paginator").Parse(pagination_tmpl)

		if(err!=nil){
			return template.HTML(err.Error())
		}

	}

	var w bytes.Buffer
	err2 := tmpl.Execute(&w,p)

	if(err2!=nil){

		return template.HTML(err2.Error())
	}

	return template.HTML(w.String())

}
