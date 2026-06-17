package value

// value/ holds value objects — small immutable computed types shared across
// layers, with no persistence or transport concerns.

type CalculatedPrice struct {
	BaseFare float64
	UnitFare float64
	Total    float64
}
