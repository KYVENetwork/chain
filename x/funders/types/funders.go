package types

func (f *Funding) AddAmount(amount uint64) {
	f.Amount += amount
}

func (f *Funding) SubtractAmount(amount uint64) (subtracted uint64) {
	subtracted = amount
	if f.Amount < amount {
		subtracted = f.Amount
	}
	f.Amount -= subtracted
	return subtracted
}

func (f *Funding) ChargeOneBundle() (amount uint64) {
	amount = f.SubtractAmount(f.AmountPerBundle)
	f.TotalFunded += amount
	return amount
}

func (f *Funding) IsActive() (isActive bool) {
	return f.Amount > 0
}

func (f *Funding) IsInactive() (isInactive bool) {
	return !f.IsActive()
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
