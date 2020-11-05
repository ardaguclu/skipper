package tee

import (
	"github.com/ardaguclu/skipper/filters"
	teepredicate "github.com/ardaguclu/skipper/predicates/tee"
	log "github.com/sirupsen/logrus"
)

const FilterName = "teeLoopback"

type teeLoopbackSpec struct{}
type teeLoopbackFilter struct {
	teeKey string
}

func (t *teeLoopbackSpec) Name() string {
	return FilterName
}

func (t *teeLoopbackSpec) CreateFilter(args []interface{}) (filters.Filter, error) {

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}
	teeKey, _ := args[0].(string)
	if teeKey == "" {
		return nil, filters.ErrInvalidFilterParameters
	}
	return &teeLoopbackFilter{
		teeKey,
	}, nil
}

func NewTeeLoopback() filters.Spec {
	return &teeLoopbackSpec{}
}

func (f *teeLoopbackFilter) Request(ctx filters.FilterContext) {
	cc, err := ctx.Split()
	if err != nil {
		log.Errorf("teeloopback: failed to split the context request: %v", err)
		return
	}
	cc.Request().Header.Set(teepredicate.HeaderKey, f.teeKey)
	go cc.Loopback()

}

func (f *teeLoopbackFilter) Response(_ filters.FilterContext) {}
