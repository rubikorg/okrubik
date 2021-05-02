version: "3.0"
services:
	{{ .Name }}:
		build:
			context: ../cmd/{{ .Name }}
			dockerfile: {{ .Name }}.Dockerfile