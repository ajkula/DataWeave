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

  // Déplacez cet événement dans la fonction init pour qu'il soit exécuté après la création des éléments
  // Gérer les boutons bascules pour les méthodes HTTP
  document.querySelectorAll('.toggle-button').forEach(button => {
    button.addEventListener('click', function() {
      this.classList.toggle('active');
      // Ici, vous pouvez ajouter la logique pour activer ou désactiver la méthode HTTP
    });
  });
  
  // Ajoutez un gestionnaire d'événements pour les boutons qui affichent les sections cachées
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
  cardContainer.innerHTML = endpointCardTemplate(tables); // Passez les tables au modèle

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

// Ajoutez cette fonction pour créer un sélecteur pour les emplacements des paramètres
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

// Ajoutez cette fonction pour créer un sélecteur pour les types de sécurité
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

// Modifiez cette fonction pour accepter les tables et retourner le modèle avec le select intégré
function endpointCardTemplate(tables) {
  // Créez un élément select pour les tables
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
        <!-- Ajoutez d'autres champs de formulaire ici -->
      </div>
      <!-- Répétez pour les sections Responses, Security -->
    </div>
  </div>
  `;
}

export async function getTranslations(msg = null) {
  return {
    pageTitle: "REST API",
  }
};
