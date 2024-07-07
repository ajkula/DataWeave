import jsYaml from 'js-yaml';
import { GetTablesList, GenerateOpenApi } from '../../../wailsjs/go/main/App';
import './styles.css';

export const html = `
<div id="apiConfigPage">
  <h1>string:pageTitle;</h1>
  <div class="configPanel">
    <div id="tablesList"></div>
    <div id="endpointConfig"></div>
  </div>
  <div id="previewPanel"></div>
  <div class="button-container">
    <button id="generateButton" class="button">string:generateButton;</button>
    <button id="resetButton" class="button">string:resetButton;</button>
  </div>
</div>
`;

let dbMetadata;
let apiConfig = {};
let currentTable = null;

export async function init() {
  dbMetadata = await GetTablesList();
  renderInterface();
  setupButtons();
}

function renderInterface() {
  renderTablesList();
  if (currentTable) {
    renderTableEndpoints(currentTable);
  } else {
    document.getElementById('endpointConfig').innerHTML = '';
  }
}

function renderTablesList() {
  const tablesList = document.getElementById('tablesList');
  tablesList.innerHTML = '';

  dbMetadata.forEach(table => {
    const tableElement = document.createElement('div');
    tableElement.className = 'table-item';
    if (currentTable && currentTable.tableName === table.tableName) {
      tableElement.classList.add('active');
    }
    tableElement.textContent = table.tableName;
    tableElement.onclick = () => {
      currentTable = table;
      renderInterface();
    };
    tablesList.appendChild(tableElement);
  });
}

function renderTableEndpoints(table) {
  const configPanel = document.getElementById('endpointConfig');
  configPanel.innerHTML = '';

  const endpoints = generateEndpointsForTable(table);
  endpoints.forEach(endpoint => {
    const endpointElement = createEndpointConfig(table, endpoint);
    configPanel.appendChild(endpointElement);
  });
}

function generateEndpointsForTable(table) {
  const tableName = table.tableName.toLowerCase();
  const endpoints = [
    { method: 'GET', path: `/${tableName}`, description: `List all ${table.tableName}`, hasFilters: true },
    { method: 'POST', path: `/${tableName}`, description: `Create a new ${table.tableName}` },
    { method: 'GET', path: `/${tableName}/{id}`, description: `Get a specific ${table.tableName}` },
    { method: 'PUT', path: `/${tableName}/{id}`, description: `Update a ${table.tableName}` },
    { method: 'PATCH', path: `/${tableName}/{id}`, description: `Partially update a ${table.tableName}` },
    { method: 'DELETE', path: `/${tableName}/{id}`, description: `Delete a ${table.tableName}` },
  ];

  table.relationships.forEach(relation => {
    const relatedTableName = relation.RelatedTableName.toLowerCase();
    endpoints.push({
      method: 'GET',
      path: `/${tableName}/{id}/${relatedTableName}`,
      description: `Get ${relation.RelatedTableName} related to ${table.tableName}`
    });

    if (relation.RelationType === 'oneToMany' || relation.RelationType === 'manyToMany') {
      endpoints.push({
        method: 'POST',
        path: `/${tableName}/{id}/${relatedTableName}`,
        description: `Add a ${relation.RelatedTableName} to ${table.tableName}`
      });
      endpoints.push({
        method: 'DELETE',
        path: `/${tableName}/{id}/${relatedTableName}/{relatedId}`,
        description: `Remove a ${relation.RelatedTableName} from ${table.tableName}`
      });
    }
  });

  return endpoints;
}

function createEndpointConfig(table, endpoint) {
  const endpointElement = document.createElement('div');
  endpointElement.className = 'endpoint-config';

  const endpointTitle = document.createElement('h3');
  endpointTitle.textContent = `${endpoint.method} ${endpoint.path}`;
  endpointElement.appendChild(endpointTitle);

  const descriptionElement = document.createElement('p');
  descriptionElement.textContent = endpoint.description;
  endpointElement.appendChild(descriptionElement);

  // Checkbox to include/exclude the endpoint
  const includeContainer = document.createElement('div');
  includeContainer.className = 'endpoint-include checkbox-container';
  const includeCheckbox = document.createElement('input');
  includeCheckbox.type = 'checkbox';
  includeCheckbox.id = `include-${endpoint.method}-${endpoint.path}`;
  includeCheckbox.checked = apiConfig[table.tableName.toLowerCase()]?.[endpoint.path]?.[endpoint.method]?.included || false;
  includeCheckbox.onchange = () => {
    updateApiConfig(table, endpoint, 'included', includeCheckbox.checked);
    updateFilterCheckboxes(endpointElement, includeCheckbox.checked);
  };
  const includeLabel = document.createElement('label');
  includeLabel.htmlFor = includeCheckbox.id;
  includeLabel.textContent = 'Include this endpoint';
  includeContainer.appendChild(includeCheckbox);
  includeContainer.appendChild(includeLabel);
  endpointElement.appendChild(includeContainer);

  // Security configuration
  const securityContainer = document.createElement('div');
  securityContainer.className = 'security-select';
  const securityLabel = document.createElement('label');
  securityLabel.textContent = 'Security: ';
  const securitySelect = document.createElement('select');
  securitySelect.innerHTML = `
    <option value="none">No security</option>
    <option value="basic">Basic Auth</option>
    <option value="bearer">Bearer Token</option>
  `;
  securitySelect.value = apiConfig[table.tableName.toLowerCase()]?.[endpoint.path]?.[endpoint.method]?.security || 'none';
  securitySelect.onchange = () => updateApiConfig(table, endpoint, 'security', securitySelect.value);
  securityContainer.appendChild(securityLabel);
  securityContainer.appendChild(securitySelect);
  endpointElement.appendChild(securityContainer);

  // Filter configuration for GET endpoints
  if (endpoint.hasFilters) {
    const filterConfig = document.createElement('div');
    filterConfig.className = 'filter-config';
    filterConfig.innerHTML = '<h4>Filterable Fields:</h4>';

    const filterFieldsContainer = document.createElement('div');
    filterFieldsContainer.className = 'filter-fields-container';

    table.columns.forEach(column => {
      const filterItem = document.createElement('div');
      filterItem.className = 'filter-item checkbox-container';
      const columnCheckbox = document.createElement('input');
      columnCheckbox.type = 'checkbox';
      columnCheckbox.id = `filter-${endpoint.method}-${endpoint.path}-${column.columnName}`;
      columnCheckbox.checked = apiConfig[table.tableName.toLowerCase()]?.[endpoint.path]?.[endpoint.method]?.filters?.[column.columnName] || false;
      columnCheckbox.onchange = () => updateApiConfig(table, endpoint, 'filters', {
        [column.columnName]: columnCheckbox.checked
      });
      const columnLabel = document.createElement('label');
      columnLabel.htmlFor = columnCheckbox.id;
      columnLabel.textContent = column.columnName;
      filterItem.appendChild(columnCheckbox);
      filterItem.appendChild(columnLabel);
      filterFieldsContainer.appendChild(filterItem);
    });

    filterConfig.appendChild(filterFieldsContainer);
    endpointElement.appendChild(filterConfig);
  }

  // Headers configuration
  const headersConfig = document.createElement('div');
  headersConfig.className = 'headers-config';
  headersConfig.innerHTML = '<h4>Headers Configuration:</h4>';

  const requestHeadersContainer = createHeadersSection('Request Headers', table, endpoint, 'requestHeaders');
  const responseHeadersContainer = createHeadersSection('Response Headers', table, endpoint, 'responseHeaders');

  headersConfig.appendChild(requestHeadersContainer);
  headersConfig.appendChild(responseHeadersContainer);
  endpointElement.appendChild(headersConfig);

  return endpointElement;
}

function updateFilterCheckboxes(endpointElement, checked) {
  const filterCheckboxes = endpointElement.querySelectorAll('.filter-config input[type="checkbox"]');
  filterCheckboxes.forEach(checkbox => {
    checkbox.checked = checked;
    checkbox.dispatchEvent(new Event('change'));
  });
}

function updateApiConfig(table, endpoint, property, value) {
  const tableName = table.tableName.toLowerCase();
  if (!apiConfig[tableName]) {
    apiConfig[tableName] = {};
  }
  if (!apiConfig[tableName][endpoint.path]) {
    apiConfig[tableName][endpoint.path] = {};
  }
  if (!apiConfig[tableName][endpoint.path][endpoint.method]) {
    apiConfig[tableName][endpoint.path][endpoint.method] = {};
  }

  const methodConfig = apiConfig[tableName][endpoint.path][endpoint.method];

  switch (property) {
    case 'filters':
      if (!methodConfig.filters) {
        methodConfig.filters = {};
      }
      Object.assign(methodConfig.filters, value);
      break;
    case 'requestHeaders':
    case 'responseHeaders':
      if (!methodConfig[property]) {
        methodConfig[property] = {};
      }
      const [headerName, isChecked] = Object.entries(value)[0];
      methodConfig[property][headerName] = isChecked;
      break;
    default:
      methodConfig[property] = value;
  }

  console.log('Updated apiConfig:', JSON.stringify(apiConfig, null, 2));
  renderInterface();
}

function setupButtons() {
  const generateButton = document.getElementById('generateButton');
  generateButton.onclick = generateOpenAPISpec;

  const resetButton = document.getElementById('resetButton');
  resetButton.onclick = resetAllOptions;
}

function resetAllOptions() {
  apiConfig = {};
  renderInterface();
}

async function generateOpenAPISpec() {
  try {
    console.log({ apiConfig });
    const spec = await GenerateOpenApi(apiConfig);
    const previewPanel = document.getElementById('previewPanel');
    previewPanel.innerHTML = '<pre>' + jsYaml.dump(jsYaml.load(spec)) + '</pre>';
  } catch (error) {
    console.error('Error generating OpenAPI spec:', error);
    // donot forget the "result" div
  }
}

function createHeadersSection(title, table, endpoint, headerType) {
  const container = document.createElement('div');
  container.className = 'headers-section';
  container.innerHTML = `<h5>${title}</h5>`;

  const headersList = getHeadersList(headerType);
  const headersContainer = document.createElement('div');
  headersContainer.className = 'headers-container';

  headersList.forEach(header => {
    const headerItem = document.createElement('div');
    headerItem.className = 'header-item checkbox-container';
    const headerCheckbox = document.createElement('input');
    headerCheckbox.type = 'checkbox';
    headerCheckbox.id = `${headerType}-${endpoint.method}-${endpoint.path}-${header}`;

    // check if fied exists
    const headerConfig = apiConfig[table.tableName.toLowerCase()]?.[endpoint.path]?.[endpoint.method]?.[headerType];
    headerCheckbox.checked = headerConfig ? headerConfig[header] || false : false;

    headerCheckbox.onchange = (event) => {
      updateApiConfig(table, endpoint, headerType, {
        [header]: event.target.checked
      });
      console.log(`Header ${header} changed to ${event.target.checked}`);
    };

    const headerLabel = document.createElement('label');
    headerLabel.htmlFor = headerCheckbox.id;
    headerLabel.textContent = header;
    headerItem.appendChild(headerCheckbox);
    headerItem.appendChild(headerLabel);
    headersContainer.appendChild(headerItem);
  });

  container.appendChild(headersContainer);
  return container;
}

function getHeadersList(headerType) {
  const requestHeaders = [
    'Accept',
    'Accept-Language',
    'Authorization',
    'Content-Type',
    'User-Agent',
    'X-Request-ID',
    'X-Correlation-ID',
    'X-Forwarded-For',
    'X-Forwarded-Host',
    'If-Modified-Since',
    'If-None-Match',
    'Cache-Control',
    'X-CSRF-Token',
    'X-API-Key',
    'X-Access-Token',
    'X-OAuth-Scopes',
    'X-Requested-With',
    'Origin',
    'Referer',
    'Cookie',
    'X-Device-ID',
    'X-App-Version'
  ];

  const responseHeaders = [
    'Content-Type',
    'ETag',
    'Last-Modified',
    'Cache-Control',
    'X-Request-ID',
    'X-Correlation-ID',
    'X-RateLimit-Limit',
    'X-RateLimit-Remaining',
    'X-RateLimit-Reset',
    'X-Content-Type-Options',
    'X-Frame-Options',
    'X-XSS-Protection',
    'Strict-Transport-Security',
    'Access-Control-Allow-Origin',
    'Access-Control-Allow-Methods',
    'Access-Control-Allow-Headers',
    'Access-Control-Expose-Headers',
    'Access-Control-Max-Age',
    'X-Powered-By',
    'X-Version',
    'Location',
    'Vary',
    'WWW-Authenticate',
    'X-CSRF-Token',
    'Set-Cookie',
    'X-Content-Security-Policy',
    'Referrer-Policy',
    'X-Download-Options',
    'X-Permitted-Cross-Domain-Policies',
    'Expect-CT',
    'Feature-Policy',
    'Public-Key-Pins',
    'Timing-Allow-Origin',
    'X-DNS-Prefetch-Control'
  ];

  return headerType === 'requestHeaders' ? requestHeaders : responseHeaders;
}

export async function getTranslations() {
  return {
    pageTitle: 'Configuration de l\'API OpenAPI',
    generateButton: 'Générer la spécification OpenAPI',
    resetButton: 'Réinitialiser toutes les options',
  };
}
