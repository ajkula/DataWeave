import { GenerateOpenApi } from '../../../wailsjs/go/main/App';
import "./styles.css";

export const html = `
<div id="apiColContainer">
  <h1>string:pageTitle;</h1>
  <div id="result" class="result"></div>
</div>
`;

export async function init() {
  const openapi = await GenerateOpenApi();
  console.log({openapi});
  const resultDiv = document.getElementById('result');
  resultDiv.innerHTML = `<pre>${CSS.escape(openapi)}</pre>`;
};

export async function getTranslations(msg = null) {
  return {
    pageTitle: "REST API",
    addCard: "Ajouter un EndPoint",
    submitBtn: "Créer l'API",
  }
};
