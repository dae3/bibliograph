import cytoscape from 'cytoscape';
import React from 'react';
import { useEffect, useState, useRef } from 'react';

export default function Graph({ books, nodeSelected }) {
  const graphdiv = useRef(null)

  const nodes = books.map(d => ({ group: 'nodes', data: d }))
  const edges = books
    .filter(d => Object.hasOwn(d, 'references'))
    .filter(d => d.references.length > 0)
    .flatMap(src => src.references.map(
      tgt => ({group: 'edges', data: {source: src.id, target: tgt}}))
    )

  useEffect(() => {
    const graph = cytoscape({
      container: graphdiv.current,
      elements: nodes.concat(edges),
      layout: { name: 'breadthfirst' },
      // userZoomingEnabled: false,
      // userPanningEnabled: false,
      style: [
        {
          selector: 'edge',
          style: {
            'curve-style': 'bezier',
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

    graph.on('click', 'node', (e) => { nodeSelected(e.target.id()) })
    graph.on('tap', 'node', (e) => { nodeSelected(e.target.id()) })
  }, [ graphdiv, books ])

  return(<div id="graph" className="min-h-lvh border" ref={graphdiv} />)
}
