import { ConfigureGorm } from '../../../wailsjs/go/main/App';
import { pagesKeys, loadPage } from '../../main';
import './styles.css';

export const html = `<div id="connectionPage">
<h1>string:pageName;</h1>

<div id="result" class="result"></div>
<form id="connectionForm">
    <label for="host">Hôte :</label>
    <input type="text" id="host" name="host">

    <label for="port">Port :</label>
    <input type="text" id="port" name="port">

    <label for="database">Type de la base de données :</label>
    <select name="database" id="database">
        <option value="postgres">PostgreSQL</option>
        <option value="mysql">MySQL</option>
        <option value="sqlite">SQLite</option>
        <option value="sqlserver">SQL Server</option>
    </select>

    <label for="username">Nom d'utilisateur :</label>
    <input type="text" id="username" name="username">

    <label for="password">Mot de passe :</label>
    <input type="password" id="password" name="password">

    <label for="dbname">Nom de la base de données :</label>
    <input type="text" id="dbname" name="dbname">

    <button type="submit" class="btn-green">Connexion à la DB</button>
</form>
</div>`

export async function init() {
    document.querySelector('#connectionForm').addEventListener('submit', function(event) {
        event.preventDefault();

        const host = document.getElementById("host").value;
        const port = document.getElementById("port").value;
        const database = document.getElementById("database").value;
        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;
        const dbname = document.getElementById("dbname").value;
        const resultDiv = document.getElementById("result");

        try {
            ConfigureGorm(database, host, port, dbname, username, password)
                .then(() => {
                   loadPage(pagesKeys.graph);
                })
                .catch((err) => {
                    resultDiv.style.display = 'flex';
                    resultDiv.innerText = err;
                });
        } catch (err) {
            console.error(err.message);
            resultDiv.style.display = 'flex';
            resultDiv.innerText = err;
        }
    });
}

// As of now, simply returns the key/val pairs
export async function getTranslations(msg = null) {
  // Here FetchTranslations will be called in near future
  return {
    pageName: 'Connexion à la base de données',
    // Add more trads HERE
    // You will also need to place it on the html like: string:your_var;
  };
}
