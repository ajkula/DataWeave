import { GraphTransform } from '../../../wailsjs/go/main/App';
import * as d3 from 'd3';
import './styles.css';
import { mergeArraysSafe } from '../../utils/utils';

export const html = `
<div id="graphPage">
  <h1>string:pageName;</h1>

  <div id="result" class="result"></div>
  <svg id="svg"></svg>
</div>
`

export async function init(data) {
  data['observer'].disconnect();
  const res = await GraphTransform();
  const graph = JSON.parse(res);

  console.log({res})
  console.log({graph})

  graph.edges.forEach(edge => {
    if (!graph.nodes.find(node => node.data.name === edge.data.source || node.data.name === edge.data.target)) {
      console.error('Missing node for edge:', edge);
    }
  });

  graph.edges.forEach(edge => {
    edge.source = graph.nodes.find(node => node.data.name === edge.data.source).data.name;
    edge.target = graph.nodes.find(node => node.data.name === edge.data.target).data.name;
  });

  const svgElement = d3.select("#svg").attr('width', '100%').attr('height', '100%').node();
  const container = svgElement.parentNode;

  // getBoundingClientRect pour les dimensions réelles
  const { width, height } = container.getBoundingClientRect();

  // Créer l'élément SVG
  const svg = d3.select(svgElement)
  .call(d3.zoom().on("zoom", ({ transform }) => svg.attr("transform", transform)))
  .append("g");

  // Créer la simulation de force
  const simulation = d3.forceSimulation(graph.nodes)
  .force("link", d3.forceLink(graph.edges).id(d => d.data.name)
    .distance(400)) // Augmenter la distance des liens
  .force("charge", d3.forceManyBody().strength(-200)) // Augmenter la répulsion
  .force("center", d3.forceCenter(width / 2, height / 2))
  .force("collide", d3.forceCollide(30));

  // Créer les liens
  const link = svg.selectAll(".link")
    .data(graph.edges)
    .enter().append("line")
    .attr("class", "link")
    .attr("marker-end", "url(#end)");

  // Créer les nœuds
  const node = svg.selectAll(".node")
    .data(graph.nodes)
    .enter().append("g")
    .attr("class", "node");

  // Ajouter des rectangles pour représenter les cartes
  node.append("rect")
  //  .attr("width", calculateNodeWidth)
    .attr("height", d => 20 + d.data.columns.length * 15 + 10)
    .attr("fill", "#fff")
    .attr("stroke", "#999");

  // Ajouter des titres aux cartes
  node.append("text")
  //  .attr("x", 60)
    .attr("y", 15)
    .attr("text-anchor", "middle")
    .text(d => d.data.name)
    .style("font-weight", "bold");

  // Ajouter les noms des colonnes
node.each(function(d) {
  const g = d3.select(this);
  d.data.columns.forEach((col, i) => {
    g.append("text")
      .attr("x", 10) // Position initiale pour les colonnes
      .attr("y", 35 + i * 15)
      .text(col.columnName)
      .style("font-size", "12px")
      .attr("text-anchor", "start"); // Aligner le texte à gauche
  });
});

  node.each(function(d) {
    const texts = d3.select(this).selectAll("text"); // Sélectionner tous les textes du nœud
    let maxWidth = 0;
  
    texts.each(function() {
      const textWidth = this.getComputedTextLength();
      if (textWidth > maxWidth) {
        maxWidth = textWidth;
      }
    });
  
    const padding = 20;
    const rectWidth = maxWidth + padding; // Ajouter le padding pour la largeur
  
    // Mettre à jour le rectangle pour avoir la largeur calculée
    d3.select(this).select("rect")
      .attr("width", rectWidth)
      .attr("x", -rectWidth / 2); // Centre le rectangle autour de l'origine du groupe
  
    // Mettre à jour la position x du texte du nom de la table pour le centrer
    d3.select(this).select("text").attr("x", 0);
  
    // Mettre à jour la position x des textes des colonnes
    d3.select(this).selectAll("text").each(function(d, i) {
      if (i > 0) { // Ignorer le premier texte (nom de la table)
        d3.select(this).attr("x", -rectWidth / 2 + 10); // Aligner à la bordure gauche du rectangle
      }
    });
  });

  function dragNode(simulation) {
    function dragstarted(event, d) {
      if (!event.active) simulation.alphaTarget(0.3).restart();
      d.fx = d.x;
      d.fy = d.y;
    }

    function dragged(event, d) {
      d.fx = event.x;
      d.fy = event.y;
    }

    function dragended(event, d) {
      if (!event.active) simulation.alphaTarget(0);
      d.fx = null;
      d.fy = null;
    }

    return d3.drag()
      .on("start", dragstarted)
      .on("drag", dragged)
      .on("end", dragended);
  }

  // Fonction de zoom pour le SVG
  function zoomed(event) {
    svg.attr("transform", event.transform);
  }

  // Appliquer le zoom au SVG
  svg.call(d3.zoom().on("zoom", zoomed));

  // Appliquer le drag aux nœuds
  node.call(dragNode(simulation));

  // Définir les marqueurs pour les flèches
  svg.append("defs").selectAll("marker")
    .data(["end"])      // Nom du marqueur
    .enter().append("marker")
    .attr("id", String)
    .attr("viewBox", "0 -5 10 10")
    .attr("refX", 25)   // Position par rapport au nœud cible
    .attr("refY", 0)
    .attr("markerWidth", 6)
    .attr("markerHeight", 6)
    .attr("orient", "auto")
    .append("path")
    .attr("d", "M0,-5L10,0L0,5")
    .attr("fill", "#ff0000");

  // Créer les liens avec les flèches
  link
    .attr("marker-end", "url(#end)"); // Utiliser le marqueur défini ci-dessus

  // Gérer la mise à jour des positions lors de la simulation
  simulation.on("tick", () => {
    link.attr("x1", d => d.source.x)
      .attr("y1", d => d.source.y)
      .attr("x2", d => d.target.x)
      .attr("y2", d => d.target.y);

    node.attr("transform", d => `translate(${d.x},${d.y})`);
  });

  // Démarrer la simulation
  simulation.restart();

  console.log(graph);
}

// As of now, simply returns the key/val pairs
export async function getTranslations(msg = null) {
  // Here FetchTranslations will be called in near future
  return {
    pageName: 'Graphe structurel',
  };
}
