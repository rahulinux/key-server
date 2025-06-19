{{- define "key-server.name" -}}
key-server
{{- end }}

{{- define "key-server.fullname" -}}
{{ .Release.Name }}-key-server
{{- end }}
