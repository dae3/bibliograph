import { customAlphabet } from 'nanoid'
import { parse, stringify } from 'yaml'
import { readFileSync, writeFileSync } from 'fs'

const BIBFILE="bib.yaml"
const OUTFILE="src/bib.json"

const nanoid = customAlphabet('1234567890abcdef', 10)
const bid = parse(readFileSync(BIBFILE, "utf8"))
  .map(b => b.hasOwnProperty('id') ? b : Object.defineProperty(b, 'id', {value: nanoid(), enumerable: true}))

writeFileSync(BIBFILE, stringify(bid), {flags: 'w', encoding: 'utf8'})
writeFileSync(OUTFILE, JSON.stringify(bid), {flags: 'w', encoding: 'utf8'})
