package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (f *Funding) GetScore(whitelist []*WhitelistCoinEntry) (score uint64) {
	// create map for easier lookup
	w := make(map[string]WhitelistCoinEntry)
	for _, entry := range whitelist {
		w[entry.CoinDenom] = *entry
	}

	for _, coin := range f.Amounts {
		score += uint64(w[coin.Denom].CoinWeight.MulInt64(coin.Amount.Int64()).TruncateInt64())
	}

	return
}

func (f *Funding) ChargeOneBundle() (payouts sdk.Coins) {
	amounts := f.Amounts
	f.Amounts = f.Amounts.Sub(f.AmountsPerBundle...)
	// we add the diff here in case amount per bundle is greater than amount
	diff := amounts.Sub(f.Amounts...)
	payouts.Add(diff...)
	f.TotalFunded.Add(diff...)
	return
}

func (f *Funding) IsActive() (isActive bool) {
	return !f.Amounts.IsZero()
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
