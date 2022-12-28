import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { MonitoringQuery, MonitoringDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, MonitoringQuery, MonitoringDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
