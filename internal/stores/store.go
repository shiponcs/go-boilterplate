package stores

// StoHolder aggregates every store behind one injectable struct, so core
// use-cases depend on a single value instead of a long parameter list. Add a
// field per new store and wire it in NewStoHolder.
type StoHolder struct {
	WidgetStore *WidgetStore
}

func NewStoHolder(widgetStore *WidgetStore) *StoHolder {
	return &StoHolder{
		WidgetStore: widgetStore,
	}
}
