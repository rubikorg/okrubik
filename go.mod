module github.com/rubikorg/okrubik

go 1.14

replace github.com/rubikorg/rubik v0.0.0 => ../ink

replace github.com/rubikorg/blocks v0.0.0 => ../blocks

require (
	github.com/AlecAivazis/survey/v2 v2.0.7
	github.com/BurntSushi/toml v0.3.1
	github.com/containerd/console v1.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/pkg/browser v0.0.0-20201207095918-0426ae3fba23
	github.com/printzero/tint v0.0.3
	github.com/radovskyb/watcher v1.0.7
	github.com/rubikorg/blocks v0.0.0
	github.com/rubikorg/rubik v0.0.0
	github.com/spf13/cobra v1.0.0
	go.mongodb.org/mongo-driver v1.4.4
	golang.org/x/tools v0.0.0-20200530233709-52effbd89c51
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v2 v2.2.8
)
