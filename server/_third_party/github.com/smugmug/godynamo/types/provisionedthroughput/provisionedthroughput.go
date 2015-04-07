package provisionedthroughput

type ProvisionedThroughput struct {
	ReadCapacityUnits  uint64
	WriteCapacityUnits uint64
}

func NewProvisionedThroughput() *ProvisionedThroughput {
	p := new(ProvisionedThroughput)
	return p
}

type ProvisionedThroughputDesc struct {
	LastIncreaseDateTime   float64
	LastDecreaseDateTime   float64
	ReadCapacityUnits      uint64
	WriteCapacityUnits     uint64
	NumberOfDecreasesToday uint64
}

func NewProvisionedThroughputDesc() *ProvisionedThroughputDesc {
	p := new(ProvisionedThroughputDesc)
	return p
}
