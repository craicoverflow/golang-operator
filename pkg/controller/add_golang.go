package controller

import (
	"github.com/craicoverflow/golang-operator/pkg/controller/golang"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, golang.Add)
}
