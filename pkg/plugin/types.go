package plugin

import (
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

type metrics struct {
	Metrics []metric
}

type metric struct {
	Name       string
	Labels     map[string]string
	Type       string
	Timeseries struct {
		Timestamps   []UnixTime
		DoubleValues []float64
		Int64Values  []int64
	}
}

type metricsReq struct {
	Query        string       `json:"query"`
	FromTime     time.Time    `json:"fromTime"`
	ToTime       time.Time    `json:"toTime"`
	Downsampling downsampling `json:"downsampling"`
}

type downsampling struct {
	GridAggregation gridAggregation `json:"gridAggregation"`
	GapFilling      gapFilling      `json:"gapFilling"`
	MaxPoints       int             `json:"maxPoints,omitempty"`
	GridInterval    Milliseconds    `json:"gridInterval,omitempty"`
	Disabled        bool            `json:"disabled,omitempty"`
}

var _ encoding.TextMarshaler = gridAggregation(0)

type gridAggregation int

const (
	gaAVG = iota
	gaMAX
	gaMIN
	gaSUM
	gaLAST
	gaCOUNT
)

func (g gridAggregation) MarshalText() (text []byte, err error) {
	switch g {
	case gaAVG:
		return []byte("AVG"), nil
	case gaMAX:
		return []byte("MAX"), nil
	case gaMIN:
		return []byte("MIN"), nil
	case gaSUM:
		return []byte("SUM"), nil
	case gaLAST:
		return []byte("LAST"), nil
	case gaCOUNT:
		return []byte("COUNT"), nil
	default:
		return nil, fmt.Errorf("unknown grid aggregation value")
	}
}

var _ encoding.TextMarshaler = gapFilling(0)

type gapFilling int

const (
	gfNONE = iota
	gfNULL
	gfPREVIOUS
)

func (g gapFilling) MarshalText() (text []byte, err error) {
	switch g {
	case gfNONE:
		return []byte("NONE"), nil
	case gfNULL:
		return []byte("NULL"), nil
	case gfPREVIOUS:
		return []byte("PREVIOUS"), nil
	default:
		return nil, fmt.Errorf("unknown gap filling value")
	}
}

var _ json.Marshaler = Milliseconds(0)

type Milliseconds time.Duration

func (m Milliseconds) MarshalJSON() ([]byte, error) {
	v := time.Duration(m).Milliseconds()
	return json.Marshal(v)
}

var _ json.Unmarshaler = (*UnixTime)(nil)

type UnixTime time.Time

func (t *UnixTime) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*(t) = UnixTime(time.UnixMilli(v))
	return nil
}
