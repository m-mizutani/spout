package memdb

import "github.com/m-mizutani/spout/pkg/model"

type MemDB struct {
	logs []*model.Log
}

func New() *MemDB {
	return &MemDB{}
}

func (x *MemDB) Put(ctx *model.Context, logs ...*model.Log) error {
	x.logs = append(x.logs, logs...)
	return nil
}

func (x *MemDB) Get(ctx *model.Context, opt *model.RepositoryGetOption) ([]*model.Log, error) {
	count := uint64(len(x.logs))

	head := opt.Offset
	if count <= head {
		return nil, nil
	}

	var results []*model.Log
	for i := head; len(results) < int(opt.Limit) && i < count; i++ {
		if opt.Filter != nil {
			logs := opt.Filter(x.logs[i])
			results = append(results, logs...)
		} else {
			results = append(results, x.logs[i])
		}
	}

	return results, nil
}
