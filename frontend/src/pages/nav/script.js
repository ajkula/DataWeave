import { pagesKeys, loadPage } from '../../main';
import './styles.css';

export const html = `<div></div>
`

export const init = async (data) => {
  const { pageName } = data;
  const navDiv = document.getElementById('nav');
  const items = [
    {
      text: 'Connexion',
      page: pagesKeys.connection,
      name: pagesKeys.connection,
    },
    {
      text: 'Graph visualizer',
      page: pagesKeys.graph,
      name: pagesKeys.graph,
    },
    {
      text: 'Referencial integrity',
      page: pagesKeys.integrity,
      name: pagesKeys.integrity,
    },
    {
      text: 'SQL Generator',
      page: pagesKeys.graph,
      name: null,
    },
    {
      text: 'REST API',
      page: pagesKeys.graph,
      name: null,
    },
  ];

  navDiv.innerHTML = '';
  items.forEach(item => {
    const div = document.createElement('div');
    div.className = 'item';
    div.textContent = item.text
    item.name === pageName ? div.style.backgroundColor = 'rgb(20, 35, 206)' : 'rgb(8, 16, 100)';
    div.onclick = () => loadPage(item.page);
    navDiv.appendChild(div);
  });
}