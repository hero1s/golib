
package admin

var indexTpl = `
{{define "content"}}
<h1>Gin Admin Dashboard</h1>
<p>
Gin 进程内监控管理后台说明:
	我不想告诉你
</p>
{{.Content}}
{{end}}`

var profillingTpl = `
{{define "content"}}
<h1>{{.Title}}</h1>
<pre id="content">
<div>{{.Content}}</div>
</pre>
{{end}}`

var defaultScriptsTpl = ``

var gcAjaxTpl = `
{{define "scripts"}}
<script type="text/javascript">
	var app = app || {};
(function() {
	app.$el = $('#content');
	app.getGc = function() {
		var that = this;
		$.ajax("/prof?command=gc%20summary&format=json").done(function(data) {
			that.$el.append($('<p>' + data.Content + '</p>'));
		});
	};
	$(document).ready(function() {
		setInterval(function() {
			app.getGc();
		}, 3000);
	});
})();
</script>
{{end}}
`

var qpsTpl = `{{define "content"}}
<h1>Requests statistics</h1>
<table class="table table-striped table-hover ">
	<thead>
	<tr>
	{{range .Content.Fields}}
		<th>
		{{.}}
		</th>
	{{end}}
	</tr>
	</thead>

	<tbody>
	{{range $i, $elem := .Content.Data}}

	<tr>
	    <td>{{index $elem 0}}</td>
	    <td>{{index $elem 1}}</td>
	    <td>{{index $elem 2}}</td>
	    <td data-order="{{index $elem 3}}">{{index $elem 4}}</td>
	    <td data-order="{{index $elem 5}}">{{index $elem 6}}</td>
	    <td data-order="{{index $elem 7}}">{{index $elem 8}}</td>
	    <td data-order="{{index $elem 9}}">{{index $elem 10}}</td>
	</tr>
	{{end}}
	</tbody>

</table>
{{end}}`

var tasksTpl = `{{define "content"}}

<h1>{{.Title}}</h1>

{{if .Message }}
{{ $messageType := index .Message 0}}
<p class="message
{{if eq "error" $messageType}}
bg-danger
{{else if eq "success" $messageType}}
bg-success
{{else}}
bg-warning
{{end}}
">
{{index .Message 1}}
</p>
{{end}}


<table class="table table-striped table-hover ">
<thead>
<tr>
{{range .Content.Fields}}
<th>
{{.}}
</th>
{{end}}
</tr>
</thead>

<tbody>
{{range $i, $slice := .Content.Data}}
<tr>
	{{range $slice}}
	<td>
	{{.}}
	</td>
	{{end}}
	<td>
	<a class="btn btn-primary btn-sm" href="/task?taskname={{index $slice 0}}">Run</a>
	</td>
</tr>
{{end}}
</tbody>
</table>

{{end}}`

// The base dashboardTpl
var dashboardTpl = `
<!DOCTYPE html>
<html lang="en">
<head>
<!-- Meta, title, CSS, favicons, etc. -->
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">

<title>

Welcome to Gin Admin Dashboard

</title>

<link href="//maxcdn.bootstrapcdn.com/bootstrap/3.2.0/css/bootstrap.min.css" rel="stylesheet">
<link href="//cdn.datatables.net/plug-ins/725b2a2115b/integration/bootstrap/3/dataTables.bootstrap.css" rel="stylesheet">

<style type="text/css">
ul.nav li.dropdown:hover > ul.dropdown-menu {
	display: block;    
}
#logo {
	width: 102px;
	height: 32px;
	margin-top: 5px;
}
.message {
	padding: 15px;
}
</style>

</head>
<body>

<header class="navbar navbar-default navbar-static-top bs-docs-nav" id="top" role="banner">
<div class="container">
<div class="navbar-header">
<button class="navbar-toggle" type="button" data-toggle="collapse" data-target=".bs-navbar-collapse">
<span class="sr-only">Toggle navigation</span>
<span class="icon-bar"></span>
<span class="icon-bar"></span>
<span class="icon-bar"></span>
</button>

<a href="/">
<img id="logo" src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=38572500,1698919762&fm=26&gp=0.jpg"/>
</a>

</div>
<nav class="collapse navbar-collapse bs-navbar-collapse" role="navigation">
<ul class="nav navbar-nav">
<li>
<a href="/qps">
Requests statistics
</a>
</li>
<li>

<li class="dropdown">
<a href="#" class="dropdown-toggle disabled" data-toggle="dropdown">Performance profiling<span class="caret"></span></a>
<ul class="dropdown-menu" role="menu">

<li><a href="/prof?command=lookup goroutine">lookup goroutine</a></li>
<li><a href="/prof?command=lookup heap">lookup heap</a></li>
<li><a href="/prof?command=lookup threadcreate">lookup threadcreate</a></li>
<li><a href="/prof?command=lookup block">lookup block</a></li>
<li><a href="/prof?command=get cpuprof">get cpuprof</a></li>
<li><a href="/prof?command=get memprof">get memprof</a></li>
<li><a href="/prof?command=gc summary">gc summary</a></li>

</ul>
</li>

<li>
<a href="/task" class="dropdown-toggle disabled" data-toggle="dropdown">Tasks</a>
</li>

</ul>
</nav>
</div>
</header>

<div class="container">
{{template "content" .}}
</div>

<script src="//code.jquery.com/jquery-1.11.1.min.js"></script>
<script src="//maxcdn.bootstrapcdn.com/bootstrap/3.2.0/js/bootstrap.min.js"></script>
<script src="//cdn.datatables.net/1.10.2/js/jquery.dataTables.min.js"></script>
<script src="//cdn.datatables.net/plug-ins/725b2a2115b/integration/bootstrap/3/dataTables.bootstrap.js
"></script>

<script type="text/javascript">
$(document).ready(function() {
    $('.table').dataTable();
});
</script>
{{template "scripts" .}}
</body>
</html>
`
