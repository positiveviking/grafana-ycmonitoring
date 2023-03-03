import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MonitoringQuery extends DataQuery {
  folderId?: string;
  aggregation?: string;
  alias?: string;
  queryText: string;
}

export const defaultQuery: Partial<MonitoringQuery> = {
  aggregation: "AVG",
  queryText: "",
};

/**
 * These are options configured for each DataSource instance
 */
export interface MonitoringDataSourceOptions extends DataSourceJsonData {
  apiEndpoint: string;
  monitoringEndpoint: string;
  folderId: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MonitoringSecureJsonData {
  apiKeyJson?: string;
}


export const defaultSourceOptions: Partial<MonitoringDataSourceOptions> = {
  apiEndpoint: "api.cloud.yandex.net:443",
}
