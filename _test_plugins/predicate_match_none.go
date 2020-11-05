package main

import (
	"net/http"

	"github.com/ardaguclu/skipper/routing"
)

type noneSpec struct{}

func InitPredicate(opts []string) (routing.PredicateSpec, error) {
	return noneSpec{}, nil
}

func (s noneSpec) Name() string {
	return "None"
}
func (s noneSpec) Create(config []interface{}) (routing.Predicate, error) {
	return nonePredicate{}, nil
}

type nonePredicate struct{}

func (p nonePredicate) Match(*http.Request) bool {
	return false
}
