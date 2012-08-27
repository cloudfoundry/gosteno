package steno

//TODO: This templates should be shipped as standalone template
//			files. However, it seems that go doesn't support package
//			resource fils as binary library directly. Figure out how
//			to distribte package with assets file such as *.html, *.css

const index_page_template = ` 
<!doctype html>

<html>
	<body>
		<fieldset>
			<legend>Information about all loggers</legend>
			<dl>
				{{range $k, $v := .LoggersInfo}}
				<dt>{{$k}}</dt>
				<dd>- {{$v}}</dd>
				{{end}}
			</dl>
		</fieldset>

		<form action="loggers" method="post" name="change_level">
			<fieldset>
				<legend>Change level here</legend>
				regexp: <input type="text" name="regexp"/>
				level: 
				<select name="level">
					{{range .Levels}}
					<option value="{{.}}">{{.}}</option>
					{{end}}
				</select></br>
				<input type="submit" value="Change"/>
			</fieldset>
		</form>
	</body>
</html>
`
