package types

import "math"

// Handles the funders of a pool. Functions safely add and remove funds to funders.
// If amount drops to zero the funder is automatically removed from the list.

// AddAmountToFunder adds the given amount to an existing funder.
// If the funder does not exist, a new funder is inserted.
func (m *Pool) AddAmountToFunder(funderAddress string, amount uint64) {
	for _, v := range m.Funders {
		if v.Address == funderAddress {
			m.TotalFunds += amount
			v.Amount += amount
			return
		}
	}
	if amount > 0 {
		// Funder was not found, insert new funder
		m.Funders = append(m.Funders, &Funder{
			Address: funderAddress,
			Amount:  amount,
		})
		m.TotalFunds += amount
	}
}

// SubtractAmountFromFunder subtracts the given amount form an existing funder
// If the amount is grater or equal to the funders amount, the funder is removed.
func (m *Pool) SubtractAmountFromFunder(funderAddress string, amount uint64) {
	for i := range m.Funders {
		if m.Funders[i].Address == funderAddress {
			if m.Funders[i].Amount > amount {
				m.TotalFunds -= amount
				m.Funders[i].Amount -= amount
			} else {
				m.TotalFunds -= m.Funders[i].Amount

				// Remove funder
				m.Funders[i] = m.Funders[len(m.Funders)-1]
				m.Funders = m.Funders[:len(m.Funders)-1]
			}
			return
		}
	}
}

func (m *Pool) RemoveFunder(funderAddress string) {
	m.SubtractAmountFromFunder(funderAddress, math.MaxUint64)
}

func (m *Pool) GetFunderAmount(address string) uint64 {
	for _, v := range m.Funders {
		if v.Address == address {
			return v.Amount
		}
	}
	return 0
}

func (m *Pool) GetLowestFunder() Funder {
	if len(m.Funders) == 0 {
		return Funder{}
	}

	lowestFunder := m.Funders[0]
	for _, v := range m.Funders {
		if v.Amount < lowestFunder.Amount {
			lowestFunder = v
		}
	}
	return *lowestFunder
}
