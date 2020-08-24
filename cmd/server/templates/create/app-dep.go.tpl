package app

// Dependency of routers and their controllers
type Dependency struct {
    ProjectConfig
}

var dep Dependency

// SetDep returns the initialized dependency of Rubik app
func SetDep(d Dependency) {
	dep = d
}

// GetDep returns the initialized dependency of Rubik app
func GetDep() Dependency {
	return dep
}
