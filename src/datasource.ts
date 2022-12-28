import { DataSourceInstanceSettings } from '@grafana/data';

import { MonitoringQuery, MonitoringDataSourceOptions } from './types';
import { DataSourceWithBackend } from '@grafana/runtime';

export class DataSource extends DataSourceWithBackend<MonitoringQuery, MonitoringDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MonitoringDataSourceOptions>) {
    super(instanceSettings);
  }
}
