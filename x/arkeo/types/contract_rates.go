package types

func FindRate(rates []*ContractRate, userType UserType, meterType MeterType) *ContractRate {
	if rates == nil {
		return nil
	}

	for _, rate := range rates {
		if rate.UserType == userType && rate.MeterType == meterType {
			return rate
		}
	}
	return nil
}
