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

  return endpointElement;
}

function updateFilterCheckboxes(endpointElement, checked) {
  const filterCheckboxes = endpointElement.querySelectorAll('.filter-config input[type="checkbox"]');
  filterCheckboxes.forEach(checkbox => {
    checkbox.checked = checked;
    checkbox.dispatchEvent(new Event('change'));
  });
}

function getApiConfigValue(table, endpoint, property, subProperty = null) {
  const tableName = table.tableName.toLowerCase();
  const config = apiConfig[tableName]?.[endpoint.path]?.[endpoint.method];
  
  if (!config) return null;
  
  if (subProperty && property === 'filters') {
    return config.filters?.[subProperty];
  }
  
  return config[property];
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

  if (property === 'filters') {
    if (!apiConfig[tableName][endpoint.path][endpoint.method].filters) {
      apiConfig[tableName][endpoint.path][endpoint.method].filters = {};
    }
    Object.assign(apiConfig[tableName][endpoint.path][endpoint.method].filters, value);
  } else {
    apiConfig[tableName][endpoint.path][endpoint.method][property] = value;
  }
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
    console.log({apiConfig});
    const spec = await GenerateOpenApi(apiConfig);
    const previewPanel = document.getElementById('previewPanel');
    previewPanel.innerHTML = '<pre>' + jsYaml.dump(jsYaml.load(spec)) + '</pre>';
  } catch (error) {
    console.error('Error generating OpenAPI spec:', error);
    // Afficher un message d'erreur à l'utilisateur
  }
}

export async function getTranslations() {
  return {
    pageTitle: 'Configuration de l\'API OpenAPI',
    generateButton: 'Générer la spécification OpenAPI',
    resetButton: 'Réinitialiser toutes les options',
  };
}
