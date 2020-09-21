package service

import (
	"bytes"
	"net/http"
	"strings"
)
import "html/template"

func (svc *Service) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=UTF-8")

	// TODO: buffer reuse
	var b bytes.Buffer
	data := struct {
		IP string
		Protocols []string
	} {
		IP: GetRemoteIP(r),
		Protocols: svc.activeProtocols,
	}
	err := svc.homeTmpl.Execute(&b, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		logger.Error(err.Error())
		return
	}
	w.Write(b.Bytes())
}

func (svc *Service) LoadHomePage() error {
	svc.homeTmpl = template.Must(template.New("home").Parse(ShittyPack(homeTemplate)))
	return nil
}

func ShittyPack(s string) string {
	lines := strings.Split(s, "\n")
	var output string
	//lastChar := uint8('>')
	for _, line := range lines {
		line := strings.TrimSpace(line)
		k := len(line)
		if k == 0 {
			continue
		}
		output += line
		//if lastChar == '>' && line[0] == '<' {
		//	output += line
		//} else if lastChar == ';' {
		//	output += line
		//} else {
		//	output += "\n" + line
		//}
		//lastChar = line[k-1]
	}
	return output
}

const homeTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Abreka: U UP?</title>
  <meta name="description" content="Test whether a proto:host:port is reachable.">
  <meta name="twitter:card" content="summary">
  <meta name="twitter:site" content="@generativist">
  <meta name="twitter:title" content="U UP?">
  <meta name="twitter:description" content="Test a port at your IP address is accessible via HTTP, HTTPS, UDP, TCP, or QUIC.">
  <style type="text/css">
body{font-family:Arial,Verdana,sans-sherif}
#uup-form{margin: 0 auto; margin-top:100px; width: 600px; text-align:center}
#responses{margin: 0 auto; margin-top:100px; margin-bottom:100px}
#responses > tr > td { padding: 1ex }
#ip-addr{margin:1ex;padding:5px;border:1px solid #0000aa;background-color:aliceblue;}
#tip-jar{text-align: center}
.input, .submit {display: inline-block; padding: 10px 15px; font-size: 20px; border-radius: 0}
.input { border: 1px solid lightgray}
.submit { background-color: lightgreen; border: 1px solid transparent; }
.success { background-color: lightgreen; text-align: center}
.failure { background-color: #ff3333; text-align: center}
.loading { background-color: #cccc33; text-align: center}
.delete { cursor: pointer}
// =============================================================
// The following github corner is from Tim Holman
// https://github.com/tholman/github-corners
// =============================================================
.github-corner:hover .octo-arm {animation: octocat-wave 560ms ease-in-out}
@keyframes octocat-wave {
    0% { transform: rotate(0deg) }
	20% { transform: rotate(-25deg) }
    40% { transform: rotate(10deg) }
    60% { transform: rotate(-25deg) }
	80% { transform: rotate(10deg) }
    100% { transform: rotate(0deg) }
}
@media (max-width: 500px) {
    .github-corner:hover .octo-arm { animation: none; }
    .github-corner .octo-arm { animation: octocat-wave 560ms ease-in-out; }
}
  </style>
</head style="text-align: center">
<body>
<div class="version">
	  <div class="demo version-section"><a href="https://github.com/abreka/uup" class="github-corner" aria-label="View source on GitHub">
		  <svg width="80" height="80" viewBox="0 0 250 250" style="fill:#151513; color:#fff; position: absolute; top: 0; border: 0; right: 0;" aria-hidden="true">
			<path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path>
			<path d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2" fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path>
			<path d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z" fill="currentColor" class="octo-body"></path>
		  </svg></a>
	  </div>
</div>
<form onSubmit="return uUp()">
<div id="uup-form" >
	<h1><a href="https://abreka.com">Abreka</a>: U Up?</h1>
	<div id="ip-addr">Your IP address is {{.IP}}</div>
	<p>Test whether you are reachable on a particular protocol port.</p>
	<select name="proto" id="uup-proto" class="input">
	{{range .Protocols}}
		<option value="{{ . }}">{{ . }} </option>
	{{end}}
	</select>
	<input type="text" name="port" id="uup-port" placeholder="port" class="input">
	<input type="submit" value="U Up?" class="submit">
</div>
</form>
<table id="responses">
</table>
<script>
function deleter(el) { 
	el.parentElement.remove();
}
function uUp(f) {
	const addr = {{.IP}};
    const proto = document.getElementById('uup-proto').value;
    const port = document.getElementById('uup-port').value;
    const responses = document.getElementById('responses');
	if(!port || !proto) {
		return false;
	}
	let el = document.createElement('tr');
	var uri = proto + "://" + addr + ":" + port;

	el.innerHTML = "<td>" + uri + '</td><td class="loading">...</td>';
	responses.prepend(el);

	let count = 0;
	let id = setInterval(function() {
		el.innerHTML = "<td>" + uri + '</td><td class="loading">' + count + '</td><td></td><td></td>';
		count += 1;
	}, 1000);
	fetch('uup/' + proto + "/" + port, {method: 'post'})
		.then(response => response.json())
		.then(data => {
			clearInterval(id);
			if(data.err) {
				el.innerHTML = "<td>" + uri + '</td><td class="failure">' + data.err + "</td><td>" + (new Date().toLocaleTimeString()) + "</td><td>🗑️</td>";
			} else {
				let c = (data.success) ? "success" : "failure";
				let v = (data.success) ? "up" : "down";
				el.innerHTML = "<td>" + uri + '</td><td class="' + c + '">' + v + "</td><td>" + (new Date().toLocaleTimeString()) + '</td><td onClick="deleter(this)" class="delete">🗑️</td>';
			}
		});
	return false;
}
</script>

<div id="tip-jar">
	Tip Me Ether: <br/> <img src="data:image/png;base64,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"/>
	<br/>
	0x91422ee7C97Fe488375CA49DDC08659090b1e0BA
	<br/><br/>
	Github Repo: <a href="https://github.com/abreka/uup">abreka/uup</a>
</div>
</body>
</html>
`


