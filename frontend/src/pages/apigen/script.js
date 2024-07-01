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
  <button id="generateButton">string:generateButton;</button>
</div>
`;

let dbMetadata;
let endpointConfig = {};

export async function init() {
  dbMetadata = await GetTablesList();
  
  renderTablesList();
 // setupConfigPanel();
  setupGenerateButton();
}

function renderTablesList() {
  const tablesList = document.getElementById('tablesList');
  tablesList.innerHTML = '';

  dbMetadata.forEach(table => {
    console.log(table)
    const tableElement = document.createElement('div');
    tableElement.className = 'table-item';
    tableElement.textContent = table.tableName;
    tableElement.onclick = () => showTableEndpoints(table);
    tablesList.appendChild(tableElement);
  });
}

function showTableEndpoints(table) {
  const configPanel = document.getElementById('endpointConfig');
  configPanel.innerHTML = '';

  const endpoints = generateEndpointsForTable(table);
  endpoints.forEach(endpoint => {
    const endpointElement = createEndpointConfig(endpoint);
    configPanel.appendChild(endpointElement);
  });
}

function generateEndpointsForTable(table) {
  return [
    { method: 'GET', path: `/${table}`, description: `List all ${table.name}` },
    { method: 'POST', path: `/${table}`, description: `Create a new ${table.name}` },
    { method: 'GET', path: `/${table}/{id}`, description: `Get a specific ${table.name}` },
    { method: 'PUT', path: `/${table}/{id}`, description: `Update a ${table.name}` },
    { method: 'DELETE', path: `/${table}/{id}`, description: `Delete a ${table.name}` },
    // Add more endpoints for relationships if needed
  ];
}

function createEndpointConfig(endpoint) {
  const endpointElement = document.createElement('div');
  endpointElement.className = 'endpoint-config';

  const endpointTitle = document.createElement('h3');
  endpointTitle.textContent = `${endpoint.method} ${endpoint.path}`;
  endpointElement.appendChild(endpointTitle);

  const descriptionElement = document.createElement('p');
  descriptionElement.textContent = endpoint.description;
  endpointElement.appendChild(descriptionElement);

  // Checkbox to include/exclude the endpoint
  const includeCheckbox = document.createElement('input');
  includeCheckbox.type = 'checkbox';
  includeCheckbox.checked = endpointConfig[endpoint.path]?.[endpoint.method]?.included ?? true;
  includeCheckbox.onchange = () => updateEndpointConfig(endpoint, 'included', includeCheckbox.checked);
  endpointElement.appendChild(includeCheckbox);
  endpointElement.appendChild(document.createTextNode('Include this endpoint'));

  // Security configuration
  const securitySelect = document.createElement('select');
  securitySelect.innerHTML = `
    <option value="none">No security</option>
    <option value="basic">Basic Auth</option>
    <option value="bearer">Bearer Token</option>
  `;
  securitySelect.value = endpointConfig[endpoint.path]?.[endpoint.method]?.security ?? 'none';
  securitySelect.onchange = () => updateEndpointConfig(endpoint, 'security', securitySelect.value);
  endpointElement.appendChild(securitySelect);

  return endpointElement;
}

function updateEndpointConfig(endpoint, property, value) {
  if (!endpointConfig[endpoint.path]) {
    endpointConfig[endpoint.path] = {};
  }
  if (!endpointConfig[endpoint.path][endpoint.method]) {
    endpointConfig[endpoint.path][endpoint.method] = {};
  }
  endpointConfig[endpoint.path][endpoint.method][property] = value;
}

function setupGenerateButton() {
  const generateButton = document.getElementById('generateButton');
  generateButton.onclick = () => generateOpenAPISpec();
}

async function generateOpenAPISpec() {
  try {
    console.log({endpointConfig})
    const spec = await GenerateOpenApi(endpointConfig);
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
  };
}