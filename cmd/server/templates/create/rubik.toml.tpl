name = "{{ .Name }}"
module = "{{ .ModulePath }}"

[[app]]
name = "server"
path = "./cmd/server"
watch = true
