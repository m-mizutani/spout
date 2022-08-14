package memdb

import (
	"github.com/google/uuid"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/pkg/model"
)

type MemDB struct {
	logs      []*model.Log
	iterators map[model.NextToken]*iterator
}

type iterator struct {
	idx    int
	filter func(log *model.Log) []*model.Log
}

func New() *MemDB {
	return &MemDB{
		iterators: make(map[model.NextToken]*iterator),
	}
}

func (x *MemDB) Put(ctx *model.Context, logs ...*model.Log) error {
	x.logs = append(x.logs, logs...)
	return nil
}

func (x *MemDB) Get(ctx *model.Context, input *model.RepositoryGetInput) (*model.RepositoryGetOutput, error) {
	count := len(x.logs)

	start := 0
	filter := input.Filter
	if input.Token != "" {
		if v, ok := x.iterators[input.Token]; ok {
			start = v.idx
			filter = v.filter
		} else {
			return nil, goerr.New("token is not found")
		}
	}

	var results []*model.Log
	var i int
	for i = start; len(results) < int(input.Limit) && i < count; i++ {
		if filter != nil {
			logs := filter(x.logs[i])
			results = append(results, logs...)
		} else {
			results = append(results, x.logs[i])
		}
	}

	output := &model.RepositoryGetOutput{
		Logs: results,
	}

	if i < count {
		output.NextToken = model.NextToken(uuid.NewString())
		x.iterators[output.NextToken] = &iterator{
			idx:    i,
			filter: filter,
		}
	}

	return output, nil
}
