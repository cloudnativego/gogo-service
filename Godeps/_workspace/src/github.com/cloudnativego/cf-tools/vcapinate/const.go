package main

import "text/template"

const (
	serviceTemplate = `
{
  "user-provided": [
  {{ range $i, $e := .UserProvided }}
  {{- if $i }}
    },
  {{ end -}}
    {
    "credentials": {
    {{ range $k, $v := $e.Credentials -}}
     "{{ $k }}": "{{ $v }}"{{if lastKey $k $e.Credentials}},{{ end }}
    {{ end -}}
    },
    "label": "user-provided",
    "name": "{{ $e.Name }}",
    "syslog_drain_url": "{{ $e.SyslogURL }}",
    "tags": []
  {{ end -}}
    }
  ]
}
`
)

var fns = template.FuncMap{
	"lastKey": func(key string, a map[interface{}]interface{}) bool {
		var last string
		for k := range a {
			last = k.(string)
		}
		return key != last
	},
}
