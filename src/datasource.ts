import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';

import { MonitoringQuery, MonitoringDataSourceOptions, defaultQuery } from './types';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { defaults } from 'lodash';

export class DataSource extends DataSourceWithBackend<MonitoringQuery, MonitoringDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MonitoringDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(inQuery: MonitoringQuery, scopedVars: ScopedVars): Record<string, any> {
    const tsrv = getTemplateSrv();
    const query = defaults(inQuery, defaultQuery);

    query.aggregation = tsrv.replace(query.aggregation, scopedVars);
    query.alias = tsrv.replace(query.alias, scopedVars);
    query.folderId = tsrv.replace(query.folderId, scopedVars);
    query.queryText = tsrv.replace(query.queryText, scopedVars);
    return query;
  };
}
