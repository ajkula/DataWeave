import { PerformAllVerifications, GetTablesList } from '../../../wailsjs/go/main/App';
import './styles.css'

export const html = `
<div id="integrity">
  <h1>string:pageTitle;</h1>
  <div id="result" class="result"></div>

  <div class="filterBar">
    <label for="checkTypeFilter">string:verificationsFilter; :</label>
    <select id="checkTypeFilter" class="filterInput">
      <option value="missingPrimaryKeys">string:missingPrimaryKeys;</option>
      <option value="nullableColumns">string:nullableColumns;</option>
      <option value="missingUniqueIndexes">string:missingUniqueIndexes;</option>
      <option value="foreignKeyIssues">string:foreignKeyIssues;</option>
      <option value="redundantIndexes">string:redundantIndexes;</option>
      <option value="sccs">string:sccs;</option>
      </select>

    <label for="tableFilter">string:tableFilter; :</label>
    <select id="tableFilter" class="filterInput">
      <!-- Dynamically filled -->
    </select>
  </div>

  <section id="problemsContainer" class="schemaProblemsContainer">
    <!-- Dynamically filled -->
  </section>
</div>
`

export async function init() {
  const schemaChecks = await PerformAllVerifications();
  const schemaData = JSON.parse(schemaChecks);
  console.log(schemaData)
  const tablesList = await GetTablesList();
  populateTableFilter(
    safeMap(tablesList).map(item => item.tableName),
  );
  applyFilters(schemaData);
  document.getElementById('tableFilter').addEventListener('change', () => applyFilters(schemaData));
  document.getElementById('checkTypeFilter').addEventListener('change', () => applyFilters(schemaData));
}

const safeMap = supposedArray => supposedArray ?? [];

function populateTableFilter(tables) {
  const select = document.getElementById('tableFilter');
  select.innerHTML = '<option value="All Tables">string:allTables;</option>';
  tables.forEach(table => {
    const option = document.createElement('option');
    option.value = table;
    option.textContent = table;
    select.appendChild(option);
  });
}

function applyFilters(schemaData) {
  const selectedTable = document.getElementById('tableFilter').value;
  const selectedCheckType = document.getElementById('checkTypeFilter').value;
  const problemsContainer = document.getElementById('problemsContainer');
  problemsContainer.innerHTML = '';

  const filteredProblems = getFilteredProblems(selectedTable, selectedCheckType, schemaData);
  filteredProblems.forEach(problem => {
    problemsContainer.appendChild(formatProblemToCard(problem));
  });
}

function getFilteredProblems(selectedTable, selectedCheckType, schemaData) {
  let problems = [];

  if (selectedCheckType === 'missingPrimaryKeys') {
    problems = problems.concat(safeMap(schemaData.missingPrimaryKeys));
  }
  if (selectedCheckType === 'nullableColumns') {
    problems = problems.concat(safeMap(schemaData.nullableColumns));
  }
  if (selectedCheckType === 'missingUniqueIndexes') {
    problems = problems.concat(safeMap(schemaData.missingUniqueIndexes));
  }
  if (selectedCheckType === 'foreignKeyIssues') {
    problems = problems.concat(safeMap(schemaData.foreignKeyIssues));
  }
  if (selectedCheckType === 'redundantIndexes') {
    problems = problems.concat(safeMap(schemaData.redundantIndexes));
  }
  if (selectedCheckType === 'sccs') {
    problems = problems.concat(safeMap(schemaData.sccs).reduce((acc, scc) => {
      if (scc.length > 1) {
        acc.push({ scc: scc, tableName: "SCC" });
        return acc;
      }
      return acc;
    }, []));
  }

  if (selectedTable !== 'All Tables' && selectedCheckType !== 'sccs') {
    problems = problems.filter(problem => problem.tableName === selectedTable);
  }

  console.log({ problems })
  return problems;
}

function formatProblemToCard(problem) {
  const card = document.createElement('div');
  card.className = 'problemCard';

  // Card title
  const title = document.createElement('div');
  title.className = "title";
  title.textContent = problem.tableName;
  card.appendChild(title);

  // Col name if exists
  if (problem.columnName ?? false) {
    const columnName = document.createElement('div');
    columnName.className = "column";
    columnName.textContent = problem.columnName;
    card.appendChild(columnName);
  }

  // index name...
  if (problem.indexName ?? false) {
    const indexName = document.createElement('div');
    indexName.className = "index";
    indexName.textContent = problem.indexName;
    card.appendChild(indexName);
  }

  // Related Table..
  if (problem.relatedTableName ?? false) {
    const relatedTable = document.createElement('div');
    relatedTable.className = "relatedTable";
    relatedTable.textContent = problem.relatedTableName;
    card.appendChild(relatedTable);
  }

  // Related Table..
  if (problem.scc ?? false) {
    console.log("table");
    problem.scc.forEach(table => {
      const scc = document.createElement('div');
      scc.className = "column";
      console.log(table);
      scc.textContent = table;
      card.appendChild(scc);
    })
  }

  // Description...
  const description = document.createElement('div');
  description.className = "description";
  description.textContent = problem.issueDescription;
  card.appendChild(description);

  return card;
}

// As of now, simply returns the key/val pairs
export async function getTranslations(msg = null) {
  // Here FetchTranslations will be called in near future

  return {
    pageTitle: 'Intégrité du schema',
    verificationsFilter: 'Vérification',
    tableFilter: 'Table',
    allTables: 'Toutes les tables',
    missingPrimaryKeys: 'Clés primaire manquantes',
    nullableColumns: 'NOT NULL = false',
    missingUniqueIndexes: 'Indexs manquants',
    foreignKeyIssues: 'Problèmes de foreign key',
    redundantIndexes: 'Indexs redondants',
    sccs: 'Relations circulaires',
  };
}