package types

import "time"

type HistoricalDataType int

const (
	HistoricalDataTypePrice HistoricalDataType = iota
	HistoricalDataTypeIv
	HistoricalDataTypeHv
)

type HistoricalDataParams struct {
	// Either Symbol or Option should be set
	Symbol string
	Option Option

	Which     HistoricalDataType
	BarWidth  time.Duration
	EndTime   time.Time
	Duration  time.Duration // Time to go back from EndTime
	IncludeAH bool          // Include afterhours data
}
