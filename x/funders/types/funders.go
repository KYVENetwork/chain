package types

import (
	"errors"
	"fmt"
)

func (f *Funding) AddAmount(amount uint64) {
	f.Amount += amount
}

func (f *Funding) SubtractAmount(amount uint64) {
	if f.Amount > amount {
		f.Amount -= amount
	} else {
		f.Amount = 0
	}
}

func (fh *FundingState) GetLowestFunding() (funding *Funding, err error) {
	if len(fh.ActiveFundings) == 0 {
		return nil, errors.New(fmt.Sprintf("no active fundings for pool %d", fh.PoolId))
	}

	lowestFunding := fh.ActiveFundings[0]
	for _, v := range fh.ActiveFundings {
		if v.Amount < lowestFunding.Amount {
			lowestFunding = v
		}
	}
	return lowestFunding, nil
}

// SetInactive moves a funding from active to inactive
func (fh *FundingState) SetInactive(funding *Funding) {
	// Remove funding from active fundings
	for i, v := range fh.ActiveFundings {
		if v.FunderAddress == funding.FunderAddress {
			fh.ActiveFundings[i] = fh.ActiveFundings[len(fh.ActiveFundings)-1]
			fh.ActiveFundings = fh.ActiveFundings[:len(fh.ActiveFundings)-1]
			break
		}
	}

	// Add funding to inactive fundings
	for _, v := range fh.InactiveFundings {
		if v.FunderAddress == funding.FunderAddress {
			return
		}
	}
	fh.InactiveFundings = append(fh.InactiveFundings, funding)
}

// SetActive moves a funding from inactive to active
func (fh *FundingState) SetActive(funding *Funding) {
	// Remove funding from inactive fundings
	for i, v := range fh.InactiveFundings {
		if v.FunderAddress == funding.FunderAddress {
			fh.InactiveFundings[i] = fh.InactiveFundings[len(fh.InactiveFundings)-1]
			fh.InactiveFundings = fh.InactiveFundings[:len(fh.InactiveFundings)-1]
			break
		}
	}

	// Add funding to active fundings
	for _, v := range fh.ActiveFundings {
		if v.FunderAddress == funding.FunderAddress {
			return
		}
	}
	fh.ActiveFundings = append(fh.ActiveFundings, funding)
}
