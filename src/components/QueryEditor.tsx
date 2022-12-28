import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent } from 'react';
import { Field, InlineField, InlineFieldRow, Input, Select, TextArea } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../datasource';
import { defaultQuery, MonitoringDataSourceOptions, MonitoringQuery } from '../types';


type Props = QueryEditorProps<DataSource, MonitoringQuery, MonitoringDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {

  onFolderChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, onRunQuery, query } = this.props;
    onChange({ ...query, folderId: event.target.value });
    onRunQuery();
  };

  onAggChange = (event: SelectableValue<string>) => {
    const { onChange, onRunQuery, query } = this.props;
    onChange({ ...query, aggregation: event.value || "" });
    onRunQuery();
  };

  onAliasChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, onRunQuery, query } = this.props;
    onChange({ ...query, alias: event.target.value });
    onRunQuery();
  };

  onQueryTextChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onChange, onRunQuery, query } = this.props;
    onChange({ ...query, queryText: event.target.value });
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { folderId, aggregation, alias, queryText } = query;
    const aggOptions: Array<SelectableValue<string>> = [
      { value: "AVG", label: "avg" },
      { value: "MAX", label: "max" },
      { value: "MIN", label: "min" },
      { value: "SUM", label: "sum" },
      { value: "LAST", label: "last" },
      { value: "COUNT", label: "count" },
    ];
    return (
      <div>
        <InlineFieldRow>
          <InlineField label='Folder ID'>
            <Input
              name='folderId'
              value={folderId || ""}
              onChange={this.onFolderChange}
            />
          </InlineField>
          <InlineField label='Aggregation'>
            <Select
              value={aggregation}
              options={aggOptions}
              onChange={this.onAggChange}
            />
          </InlineField>
          <InlineField label='Alias'>
            <Input
              name='alias'
              value={alias || ""}
              onChange={this.onAliasChange}
            />
          </InlineField>
        </InlineFieldRow>
        <Field label='Query' description='Plain query text. Use monitoring UI for more powerfull editing and copy result here.'
          invalid={queryText.includes("folderId")}
          error="do not use folderId in query"
        >
          <TextArea
            name='query-text'
            value={queryText}
            onChange={this.onQueryTextChange}
          />
        </Field>
      </div>
    );
  }
}
