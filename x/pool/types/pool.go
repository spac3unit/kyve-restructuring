package types

import (
	"sort"
)

func (m *Pool) AddToFunder(funderAddress string, amount uint64) {
	for _, v := range m.Funders {
		if v.Address == funderAddress {
			m.TotalFunds += amount
			v.Amount += amount
		}
	}
	sort.SliceStable(m.Funders, func(i, j int) bool {
		return m.Funders[i].Amount < m.Funders[j].Amount
	})
}

func (m *Pool) SubFromFunder(funderAddress string, amount uint64) {
	for _, v := range m.Funders {
		if v.Address == funderAddress {
			if v.Amount > amount {
				m.TotalFunds -= amount
				v.Amount -= amount
			} else if v.Amount == amount {
				m.TotalFunds -= amount
				v.Amount -= amount

				m.RemoveFunder(*v)
			}
		}
	}
	sort.SliceStable(m.Funders, func(i, j int) bool {
		return m.Funders[i].Amount < m.Funders[j].Amount
	})
}

func (m *Pool) InsertFunder(funder Funder) {

	m.Funders = append(m.Funders, &funder)
	m.TotalFunds += funder.Amount
	sort.SliceStable(m.Funders, func(i, j int) bool {
		return m.Funders[i].Amount < m.Funders[j].Amount
	})
}

func (m *Pool) RemoveFunder(funder Funder) {
	index := sort.Search(len(m.Funders), func(i int) bool {
		if m.Funders[i].Amount == funder.Amount {
			return m.Funders[i].Address >= funder.Address
		}
		return m.Funders[i].Amount >= funder.Amount
	})
	if index < len(m.Funders) {
		if m.Funders[index].Address == funder.Address {
			m.Funders = append(m.Funders[0:index], m.Funders[index+1:len(m.Funders)]...)
			m.TotalFunds -= funder.Amount
		}
	}
}

func (m *Pool) GetFunder(address string) (Funder, bool) {
	for _, v := range m.Funders {
		if v.Address == address {
			return *v, true
		}
	}
	return Funder{}, false
}

func (m *Pool) GetLowestFunder() Funder {
	if len(m.Funders) == 0 {
		return Funder{}
	}
	return *m.Funders[0]
}
