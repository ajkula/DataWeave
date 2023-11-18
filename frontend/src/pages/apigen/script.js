import { GetTablesList } from '../../../wailsjs/go/main/App';
import "./styles.css";
const HTTP_METHOD_LIST = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];
const CONTENT_TYPES = ['application/json', 'application/xml', 'text/plain', 'multipart/form-data'];
const PARAMETERS_LOCATION = ['query', 'header', 'path', 'cookie'];
const SECURITY_TYPES = ['apiKey', 'http', 'oauth2', 'openIdConnect'];
let selecteed_table = {}

export const html = `
<div id="apiColContainer">
  <h1>string:pageTitle;</h1>
  <div id="result" class="result"></div>
  <div id="apiFormContainer"></div>
</div>
`;

export async function init() {
  const tables = await GetTablesList();
  selecteed_table = tables[0];
  const formContainer = document.getElementById('apiFormContainer');
  console.log(selecteed_table);
  
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

function createLabeledSelectorFromList(arr, selectName) {
  let options = arr.map(item => `<option value="${item}">${item}</option>`).join('');
  return `
    <select class="parameter-location-select">
      <option value="">${selectName}</option>
      ${options}
    </select>
  `;
}

function endpointCardTemplate(tables) {
  let tablesOptions = tables.map(table => table.tableName);

  return `
  <div class="endpoint-card">
    <div class="endpoint-header api-card-row">
      <div class="api-card-col">
        ${createLabel('Select Table')}
        ${createLabeledSelectorFromList(tablesOptions, "Select Table")}
        ${createLabel('Select Method')}
        ${createLabeledSelectorFromList(HTTP_METHOD_LIST, "Select Method")}
      </div>
      <div class="api-card-col">
        ${createLabel("Parameters")}
        ${createLabeledSelectorFromList(PARAMETERS_LOCATION, "Parameters")}
        ${createLabel("Security Type")}
        ${createLabeledSelectorFromList(SECURITY_TYPES, "Security Type")}
      </div>
      <div class="api-card-col">
        ${createLabel("Content Type")}
        ${createLabeledSelectorFromList(CONTENT_TYPES, "Content Type")}
      </div>
    </div>
  </div>
  `;
}

export async function getTranslations(msg = null) {
  return {
    pageTitle: "REST API",
  }
};
