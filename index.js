#! /usr/local/bin/node

const monitor = require('./lib/monitor')

const args = process.argv.slice(2)

var stop = args[0].normalize('NFC')
var offset = 0

var offsetMatch = args[0].match(/in (\d+)/)
if (offsetMatch !== null && offsetMatch.length > 0) {
  offset = offsetMatch[1]
  stop = args[0].split('in')[0]
}

monitor(stop, offset)
  .then(createAlfredJSON)
  .then(JSON.stringify)
  .then(console.log)

function createAlfredJSON(items) {
  return {
    'items': items
  }
}
