package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler      = (*monitoringDatasource)(nil)
	_ backend.CheckHealthHandler    = (*monitoringDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*monitoringDatasource)(nil)
)

type monitoringDatasource struct {
	logger log.Logger

	sdk      *sdk
	folderID string
}

func NewDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var pubConfig monitoringConfig
	if err := json.Unmarshal(settings.JSONData, &pubConfig); err != nil {
		return nil, fmt.Errorf("unmarshal plugin config: %w", err)
	}

	sdk, err := newSDK(pubConfig.APIEndpoint, pubConfig.MonitoringEndpoing, settings.DecryptedSecureJSONData[apiKeyJsonInSettings])
	if err != nil {
		return nil, fmt.Errorf("yc sdk: %w", err)
	}

	return &monitoringDatasource{
		logger:   log.DefaultLogger,
		sdk:      sdk,
		folderID: pubConfig.FolderID,
	}, nil
}

func (d *monitoringDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	resp := backend.NewQueryDataResponse()

	for _, query := range req.Queries {
		var mr monitoringRequest
		if err := json.Unmarshal(query.JSON, &mr); err != nil {
			d.logger.Error("unmarshal query issue", "error", err.Error(), "ref_id", query.RefID, "query", query.JSON)
			return &backend.QueryDataResponse{}, fmt.Errorf("can not unmarshal query: %w", err)
		}

		folderID := d.folderID
		if mr.FolderID != "" {
			folderID = mr.FolderID
		}
		req := metricsReq{
			Query:    mr.QueryText,
			FromTime: query.TimeRange.From,
			ToTime:   query.TimeRange.To,
			Downsampling: downsampling{
				GridAggregation: parseAggregation(mr.Aggregation),
				MaxPoints:       int(query.MaxDataPoints),
			},
		}

		metrics, err := d.sdk.read(ctx, folderID, req)
		respD := resp.Responses[query.RefID]
		if err != nil {
			d.logger.Error("read metrics error", "error", err.Error(), "ref_id", query.RefID)
			respD.Error = err
			resp.Responses[query.RefID] = respD
			continue
		}

		for _, metric := range metrics.Metrics {
			var valuesField *data.Field
			switch {
			case len(metric.Timeseries.DoubleValues) > 0:
				valuesField = valueField(mr.Alias, metric.Name, metric.Labels, metric.Timeseries.DoubleValues)
			case len(metric.Timeseries.Int64Values) > 0:
				valuesField = valueField(mr.Alias, metric.Name, metric.Labels, metric.Timeseries.Int64Values)
			default:
				continue
			}

			timestamps := make([]time.Time, len(metric.Timeseries.Timestamps))
			for i := range timestamps {
				timestamps[i] = time.Time(metric.Timeseries.Timestamps[i])
			}

			frame := data.NewFrame(query.RefID,
				data.NewField("timestamp", nil, timestamps),
				valuesField,
			)

			frame.SetMeta(&data.FrameMeta{
				PreferredVisualization: data.VisTypeGraph,
			})
			respD.Frames = append(respD.Frames, frame)
		}
		resp.Responses[query.RefID] = respD
	}
	return resp, nil
}

func (d *monitoringDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if err := d.sdk.check(ctx); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "OK",
	}, nil
}

func (d *monitoringDatasource) Dispose() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := d.sdk.Shutdown(ctx); err != nil {
		d.logger.Error("plugin dispose error", "error", err.Error())
	}
}

func parseAggregation(in string) gridAggregation {
	switch strings.ToUpper(in) {
	default:
		return gaAVG
	case "AVG":
		return gaAVG
	case "MAX":
		return gaMAX
	case "MIN":
		return gaMIN
	case "SUM":
		return gaSUM
	case "LAST":
		return gaLAST
	case "COUNT":
		return gaCOUNT
	}
}

func valueField(alias string, name string, labels map[string]string, values interface{}) *data.Field {
	if alias != "" {
		rndr, err := mustache.Render(alias, labels)
		if err == nil {
			name = rndr
			labels = nil
		}
	}
	return data.NewField(name, labels, values)
}
