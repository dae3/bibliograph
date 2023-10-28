import './main.css';
import bib from './bib.json';
import cytoscape from 'cytoscape';

function graph() {
  // raw data shape is [ {id,author,title,references[]} ]
  const nodes = bib.map(d => ({ group: 'nodes', data: d }))
  const edges = bib
    .filter(d => Object.hasOwn(d, 'references'))
    .filter(d => d.references.length > 0)
    .flatMap(src => src.references.map(
      tgt => ({group: 'edges', data: {source: src.id, target: tgt}}))
    )
  console.log(edges)
  const graph = cytoscape({
    container: document.getElementById("graph"),
    elements: nodes.concat(edges),
    layout: { name: 'breadthfirst' },
    style: [
      {
        selector: 'edge',
        style: {
          'curve-style': 'bezier',
          'line-color': 'blue',
          'target-arrow-color': 'blue',
          'target-arrow-shape': 'triangle',
          'arrow-scale': 2
        }
      },
      {
        selector: 'node',
        style: {
          label: n => n.data('title') == null ? 'no title' : n.data('title')
        }
      }
    ]
  })
}
window.addEventListener('load', graph)
