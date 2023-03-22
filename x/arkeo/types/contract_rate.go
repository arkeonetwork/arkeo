package types

func FindRate(rates []*ContractRate, userType UserType, meterType MeterType) *ContractRate {
	for _, rate := range rates {
		if rate.UserType == userType && rate.MeterType == meterType {
			return rate
		}
	}
	return nil
}
