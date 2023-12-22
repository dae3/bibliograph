import cytoscape from 'cytoscape';
import React from 'react';
import { useEffect, useState, useRef } from 'react';

export default function Graph({ books }) {
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
  }, [ graphdiv, books ])

  return(<div id="graph" ref={graphdiv} />)
}
