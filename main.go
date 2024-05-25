<!DOCTYPE html>
<html lang="en">

<head>
	<title>SAY HI!</title>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale =1">
	<script src="https://unpkg.com/htmx.org@1.9.12"
		integrity="sha384-ujb1lZYygJmzgSwoxRggbCHcjc0rB2XoQrxeTUQyRjrOnlCoYta87iKBWq3EsdM2"
		crossorigin="anonymous"></script>
	<link href="css/style.css" rel="stylesheet">
</head>

<body>
	<button hx-post="/project/1" hx-target="main"> > Project 1 </button>
	<button hx-post="/project/2" hx-target="main"> > Project 2 </button>

	<main>
		{{template "project" .}}
	</main>
</body>

</html>

{{ block "project" .}}
{{.Project}}
{{end}}
