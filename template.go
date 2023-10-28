package table

var default_template = `
  {{$columns := .Columns }}
  <table class = "{{.Class}}">
     {{if not .Opts.HideHeader }}
     <thead>
     <tr class="bg-primary text-white">
     {{ range $k,$v := .Columns }}
		  <th   {{ if gt (len $v.Ordering) 0 }}class="sortStyle"{{end}} >
             {{$v.Render}}
          </th>
     {{end}}
     </tr>
     </thead>
     {{end}}
     {{ range $krow,$row := .Data }}
      <tr>
            {{ $rrow := (tomap $row) }}
			{{ range $k,$v := $columns }}
                 <td>
                 {{ parse $v.ID $columns $rrow }}
                 </td>
			{{end}}
      </tr>
    {{end}}
    {{ if lt (len .Data) 1 }}
    <tr>
       <td colspan="{{ len $columns }}">
			<div class="text-center" >
					 No Records Found
			</div>
       </td>
    </tr>
    {{end}}
  </table>
  {{if .Paginator}}
  {{.Paginator.Render}}
  {{end}}
`

var default_pagination = `
{{$current := .Current}}
{{$pagelink := .Url}}
{{$obj := .}}
{{ if gt .Total 0 }}
<div class="row m-1 pt-3">
<div class="col-lg-6 col-12 col-md-6" >
<ul class="pagination pagination-rounded">
    {{if .HasPrev}}
        <li class="page-item" ><a class="page-link" href="{{$obj.GetUrl 1}}">&laquo; &laquo;</a></li>
        <li class="page-item" ><a class="page-link" href="{{.PageLinkPrev}}">&laquo;</a></li>
    {{else}}
        <li class="page-item disabled"><a class="page-link" >&laquo;&laquo;</a></li>
        <li class="page-item disabled"><a class="page-link" >&laquo;</a></li>
    {{end}}
    {{range $index, $page := .Pages}}
        <li class="page-item {{if (eq $page $current) }} active{{end}}" >
            <a class="page-link" href="{{$obj.GetUrl $page}}">{{$page}}</a>
        </li>
    {{end}}
    {{if .HasNext}}
        <li class="page-item" ><a class="page-link" href="{{.PageLinkNext}}">&raquo;</a></li>
        <li class="page-item" ><a class="page-link" href="{{.PageLinkLast}}">&raquo; &raquo;</a></li>
    {{else}}
        <li class="page-item disabled"><a class="page-link" >&raquo;</a></li>
        <li class="page-item disabled"><a class="page-link">&raquo; &raquo;</a></li>
    {{end}}
</ul>
</div>
<div class="col-lg-6 col-12 col-md-6">
<div class="float-end">
<select name="perpage" onchange="window.location = this.value;">
 {{ range $k,$v := .PageList }}
	{{if (eq $obj.Perpage $v)}} 
	<option selected value="{{$obj.GetPerpage $v}}">{{$v}} Items/Page</option>
	  {{else}}
	<option value="{{$obj.GetPerpage $v}}">{{$v}} Items/Page</option>
	{{end}}
 {{ end }}
</select>
</div>
</div>
</div>
{{end}}
`
