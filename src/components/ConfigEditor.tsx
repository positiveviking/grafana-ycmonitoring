import React, { ChangeEvent, PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MonitoringDataSourceOptions, MonitoringSecureJsonData, defaultSourceOptions } from '../types';
import { defaults } from 'lodash';
import { Field, Input, SecretTextArea } from '@grafana/ui';


interface Props extends DataSourcePluginOptionsEditorProps<MonitoringDataSourceOptions> { }

interface State { }

export class ConfigEditor extends PureComponent<Props, State> {
  onEndpointChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      apiEndpoint: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onMonitoringEndpointChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      monitoringEndpoint: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onFolderChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      folderId: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onAPIKeyChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        apiKeyJson: event.target.value,
      },
    });
  };

  onResetAPIKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        apiKeyJson: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        apiKeyJson: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { secureJsonFields } = options;
    const jsonData = defaults(options.jsonData, defaultSourceOptions) as MonitoringDataSourceOptions;
    const secureJsonData = (options.secureJsonData || {}) as MonitoringSecureJsonData;

    return (
      <div>
        <Field label="API endpoint">
          <Input
            onChange={this.onEndpointChange}
            value={jsonData.apiEndpoint}
            placeholder="yandex cloud api endpoint <host>:<port>"
          />
        </Field>
        <Field label="Monitoring endpoint">
          <Input
            onChange={this.onEndpointChange}
            value={jsonData.monitoringEndpoint}
            placeholder="yandex cloud monitoring custom endpoint <host>:<port>"
          />
        </Field>
        <Field label="Folder ID">
          <Input
            onChange={this.onFolderChange}
            value={jsonData.folderId || ''}
            placeholder="folder for metrics read (required)"
          />
        </Field>

        <Field label="API Key">
          <SecretTextArea
            rows={8}
            isConfigured={(secureJsonFields && secureJsonFields.apiKeyJson) as boolean}
            value={secureJsonData.apiKeyJson || ''}
            placeholder="place full json key file content here"
            onReset={this.onResetAPIKey}
            onChange={this.onAPIKeyChange}
          />
        </Field>
      </div>
    );
  }
}
