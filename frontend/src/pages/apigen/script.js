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
  <div class="controls">
    <label for="add">string:addCard; :</label>
    <button id="add" class="addBtn">string:addCard;</button>
    <label for="submitBtn">string:submitBtn; :</label>
    <button id="submitBtn" class="submitBtn">string:submitBtn;</button>
  </div>
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
  
  initializeEventListeners();
};

function createLabel(text) {
  return `<label class="form-label" for="${text}">${text}</label>`;
}

function createEndpointCard(tables) {
  const cardContainer = document.createElement('div');
  cardContainer.innerHTML = endpointCardTemplate(tables);
  return cardContainer.firstElementChild;
}

function createLabeledSelectorFromList(arr, selectName, extraClass = "") {
  let options = arr.map(item => `<option value="${item}">${item}</option>`).join('');
  return `
    ${createLabel(selectName)}
    <select class="parameter-location-select ${extraClass}">
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
        ${createLabeledSelectorFromList(tablesOptions, "Select Table")}
        ${createLabeledSelectorFromList(HTTP_METHOD_LIST, "Select Method")}
        ${createByIdSwitch("byIdSwitch", "Ressource by ID")}
      </div>
      <div class="api-card-col">
        ${createLabeledSelectorFromList(PARAMETERS_LOCATION, "Parameters")}
        ${createLabeledSelectorFromList(SECURITY_TYPES, "Security Type")}
        ${createLabeledSelectorFromList(CONTENT_TYPES, "Content Type")}
      </div>
    </div>
  </div>
  `;
}

function createByIdSwitch(id, label) {
  return `
    ${createLabel(label)}
    <label class="switch">
        <input type="checkbox" id="${id}">
        <span class="slider round"></span>
    </label>
  `;
}

function initializeEventListeners() {
  document.getElementById('submitBtn').addEventListener('click', function() {
    console.log("CREATE !")
  });


  // Listener for "ById" switch
  const byIdSwitch = document.getElementById('byIdSwitch');
  if (byIdSwitch) {
      byIdSwitch.addEventListener('change', function(ev) {
          console.log(ev.target.checked)
          // todo: Activation logic for "ById"...
      });
  }
}

export async function getTranslations(msg = null) {
  return {
    pageTitle: "REST API",
    addCard: "Ajouter un EndPoint",
    submitBtn: "Créer l'API",
  }
};
