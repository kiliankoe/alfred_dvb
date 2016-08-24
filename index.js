#! /usr/local/bin/node

const dvb = require('dvbjs')
const moment = require('moment')

moment.locale('de')

const args = process.argv.slice(2)
const notificationOffset = 10

var stop = args[0]
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

function parseConnection(con) {

    var arg = '0 "Zu bald" "Verbindung muss für eine Benachrichtigung mehr als 10 Minuten in der Zukunft liegen."'
    if (con.arrivalTimeRelative > notificationOffset) {
        arg = `${con.arrivalTimeRelative - notificationOffset} "Auf geht's" "Die ${con.line} Richtung ${con.direction} fährt in ${notificationOffset} Minuten."`
    }

    return {
        'title': con.line + ' ' + con.direction + createArrivalTimeString(con.arrivalTimeRelative),
        'subtitle': moment().add(con.arrivalTimeRelative, 'm').format('dddd, HH:mm [Uhr]'),
        'arg': arg,
        'icon': {
            'path': `transport_icons/${con.mode.name}.png`
        }
    }
}

function monitor(stop, timeOffset = 0, numResults = 6) {
    return dvb.monitor(stop, timeOffset, numResults)
    .then((data) => {

        if (data.length === 0) {
            return [{
                'title': 'Haltestelle nicht gefunden 🤔',
                'subtitle': 'Vielleicht ein Tippfehler?'
            }]
        }

        return data.map(parseConnection)
    })
    .catch(err => {
        return [{
            'title': 'Unerwarteter Fehler 😲',
            'subtitle': err.message,
        }]
    })
}

function createArrivalTimeString(arrivalTime) {
    if (arrivalTime === 0) {
        return ' jetzt'
    } else if (arrivalTime === 1) {
        return ' in 1 Minute'
    } else {
        return ` in ${arrivalTime} Minuten`
    }
}

function createAlfredJSON(items) {
    return {
        'items': items
    }
}
