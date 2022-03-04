# go-table

Go table helps with table output

Basic Example

```
  tab := table.New()
	tab.Class = "table table-striped sortable-table"

	pg,_ := strconv.Atoi(t.Ctx.Request().URL.Query().Get("page"))
	perpage,_ := strconv.Atoi(t.Ctx.Request().URL.Query().Get("perpage"))

	col := util.UrlQuery(t.Ctx.Request(),"order.column","Account_id")
	order := util.UrlQuery(t.Ctx.Request(),"order","ASC")

	pag := tab.SetPaginator(int(totalrows),perpage,pg)

	tab.SetRequest(r) //type http.Request

	offset,limit := pag.GetOffsets()

  //the Add Column blocks will be called when iterating, expensive tasks like queries are not recommended

	tab.AddColumn("Account","LastName",func(r interface{}) (template.HTML){

		var m *models.Account

		table.ToInterface(r,&m)

		return template.HTML(fmt.Sprintf(`<a  href="/account/%d" >%s %s</a>`,m.AccountID,m.FirstName,m.LastName))

	},true)

	tab.AddColumn("Email","Email","",true) //will return just the Email value

	tab.AddColumn("Activated","Activated","{{.Activated}}",true) //executes the template
  
  tab.Data = reg
	out,err := tab.Render()
  
```
