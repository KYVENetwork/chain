package types

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

func (f *Funding) ChargeOneBundle() (amount uint64) {
	amount = f.AmountPerBundle
	if f.Amount < f.AmountPerBundle {
		amount = f.Amount
	}
	f.SubtractAmount(amount)
	return amount
}

// SetInactive removes a funding from active fundings
func (fs *FundingState) SetInactive(funding *Funding) {
	for i, funderAddress := range fs.ActiveFunderAddresses {
		if funderAddress == funding.FunderAddress {
			fs.ActiveFunderAddresses[i] = fs.ActiveFunderAddresses[len(fs.ActiveFunderAddresses)-1]
			fs.ActiveFunderAddresses = fs.ActiveFunderAddresses[:len(fs.ActiveFunderAddresses)-1]
			break
		}
	}
}

// SetActive adds a funding to active fundings
func (fs *FundingState) SetActive(funding *Funding) {
	for _, funderAddress := range fs.ActiveFunderAddresses {
		if funderAddress == funding.FunderAddress {
			return
		}
	}
	fs.ActiveFunderAddresses = append(fs.ActiveFunderAddresses, funding.FunderAddress)
}
