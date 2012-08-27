package steno

//TODO: This templates should be shipped as standalone template
//      files. However, it seems that go doesn't support package
//      resource fils as binary library directly. Figure out how
//      to distribte package with assets file such as *.html, *.css

const index_page_template = `
<!doctype html>

<html>
	<body>
    <script>
      function sendPutRequest(form, path) {
				var data = {};
				data.regexp = form.regexp.value;
				data.level =  form.level.value;
				var dataJson = JSON.stringify(data)
				url = document.URL.concat(path)
        xmlHttp = new XMLHttpRequest();
				xmlHttp.open('PUT', url, true);
				xmlHttp.setRequestHeader('Content-type','application/json; charset=utf-8');
				xmlHttp.onreadystatechange = function() {
					document.location.reload(true)
				}
        xmlHttp.send(dataJson);
				return false;
      }
    </script>
		<fieldset>
			<legend>Information about all loggers</legend>
			<dl>
				{{range $k, $v := .LoggersInfo}}
				<dt>{{$k}}</dt>
				<dd>- {{$v}}</dd>
				{{end}}
			</dl>
		</fieldset>

		<form onsubmit="return sendPutRequest(this, 'regexp')" name="change_level">
			<fieldset>
				<legend>Change level through regular expression here</legend>
				regexp: <input type="text" name="regexp"/>
				level:
				<select name="level">
					{{range .Levels}}
					<option value="{{.}}">{{.}}</option>
					{{end}}
				</select></br>
				<input type="submit" value="Change" >
			</fieldset>
		</form>
	</body>
</html>
`
