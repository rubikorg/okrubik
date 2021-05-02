name = "{{ .Name }}"
module = "{{ .ModulePath }}"

[[app]]
name = "server"
path = "./cmd/server"
watch = true
communicate = false

[x]
  [x.test]
    command = "go test -cover ./..."
    pwd = "cmd/server"