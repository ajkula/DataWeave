import { GetTablesList } from '../../../wailsjs/go/main/App';
import "./styles.css";

export const html = `
<div id="apiColContainer">
  <h1>string:pageTitle;</h1>
  <div id="result" class="result"></div>
  <div id="apiFormContainer"></div>
</div>
`;

export const HTTP_METHOD_LIST = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];

export async function init() {
  const formContainer = document.getElementById('apiFormContainer');
  const tables = await GetTablesList();
  console.log(tables);
  
  const newCard = createEndpointCard(tables);
  formContainer.appendChild(newCard);

  // HTTP methods switchs... todo
  document.querySelectorAll('.toggle-button').forEach(button => {
    button.addEventListener('click', function() {
      this.classList.toggle('active');
    });
  });
  
  // Toggle .hidden
  document.querySelectorAll('.show-section-btn').forEach(button => {
    button.addEventListener('click', function() {
      const section = document.querySelector('.' + this.dataset.section + '-section');
      section.classList.toggle('hidden');
    });
  });
};

function createLabel(text) {
  return `<label class="form-label" for="${text}">${text}</label>`;
}

function createEndpointCard(tables) {
  const cardContainer = document.createElement('div');
  cardContainer.innerHTML = endpointCardTemplate(tables);

  return cardContainer.firstElementChild;
}

function createContentTypeSelector() {
  const contentTypes = ['application/json', 'application/xml', 'text/plain', 'multipart/form-data'];
  let options = contentTypes.map(type => `<option value="${type}">${type}</option>`).join('');
  return `
    <select class="content-type-select">
      <option value="">Select Content Type</option>
      ${options}
    </select>
  `;
}

// Param source types selector
function createParameterLocationSelector() {
  const locations = ['query', 'header', 'path', 'cookie'];
  let options = locations.map(location => `<option value="${location}">${location}</option>`).join('');
  return `
    <select class="parameter-location-select">
      <option value="">Select Parameter Location</option>
      ${options}
    </select>
  `;
}

// Secu types selector
function createSecurityTypeSelector() {
  const securityTypes = ['apiKey', 'http', 'oauth2', 'openIdConnect'];
  let options = securityTypes.map(type => `<option value="${type}">${type}</option>`).join('');
  return `
    <select class="security-type-select">
      <option value="">Select Security Type</option>
      ${options}
    </select>
  `;
}

function endpointCardTemplate(tables) {
  let tablesOptions = tables.map(table => `<option value="${table.tableName}">${table.tableName}</option>`).join('');
  let methodsOptions = HTTP_METHOD_LIST.map(method => `<option value="${method}">${method}</option>`).join('');

  return `
  <div class="endpoint-card">
    <div class="endpoint-header api-card-row">
      <div class="api-card-col">
        ${createLabel('Select Table')}
        <select class="table-select">
          <option value="">Select Table</option>
          ${tablesOptions}
        </select>
        ${createLabel('Select Method')}
        <select class="method-select">
          <option value="">Select Method</option>
          ${methodsOptions}
        </select>
      </div>
      <div class="api-card-col">
        ${createLabel('Endpoint Description')}
        <input type="text" placeholder="Description" class="endpoint-description">
        ${createLabel('Endpoint Operation ID')}
        <input type="text" placeholder="Operation ID" class="endpoint-operation-id">
      </div>
    </div>
    <div class="endpoint-body">
      <button class="show-section-btn" data-section="parameters">Show Parameters</button>
      <div class="parameters-section">
        ${createParameterLocationSelector()}
        <!-- more forms prolly here -->
      </div>
      <!-- sections Responses, Security... whatever, I still gotta prep this UI -->
    </div>
  </div>
  `;
}

export async function getTranslations(msg = null) {
  return {
    pageTitle: "REST API",
  }
};
