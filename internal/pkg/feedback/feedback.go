package feedback

import (
	"aed-api-server/internal/interfaces/entities"
	"errors"
)

type ValuableFeedBack map[string]interface{}

func NewValuableFeedBack() ValuableFeedBack {
	return make(ValuableFeedBack, 1)
}

func (v ValuableFeedBack) AddPointsEventRsts(rsts []*entities.DealPointsEventRst) error {
	for _, r := range rsts {
		err := v.AddPointsEventRst(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v ValuableFeedBack) AddPointsEventRst(rst *entities.DealPointsEventRst) error {
	key := "dealPointsEventRsts"
	if r, exists := v[key]; exists {
		if dp, ok := r.([]*entities.DealPointsEventRst); ok {
			v[key] = append(dp, rst)
		} else {
			return errors.New("invalid type")
		}

		return nil
	}

	r := make([]*entities.DealPointsEventRst, 0)
	r = append(r, rst)
	v[key] = r

	return nil
}

func (v ValuableFeedBack) Put(key string, value interface{}) {
	v[key] = value
}
