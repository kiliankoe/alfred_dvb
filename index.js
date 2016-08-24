#! /usr/local/bin/node

var dvb = require('dvbjs')
var moment = require('moment')
moment.locale('de')

var args = process.argv.slice(2)

var offset = 0
var stop = args[0]

var offsetMatch = args[0].match(/in (\d+)/)
if (offsetMatch !== null && offsetMatch.length > 0) {
    offset = offsetMatch[1]
    stop = args[0].split('in')[0]
}

monitor(stop, offset).then(function (output) {
    console.log(JSON.stringify(output))
})

function monitor (stop, timeOffset = 0, numResults = 6) {
    return dvb.monitor(stop, timeOffset, numResults)
    .then(function (data) {
        var items = {'items': []}

        if (data.length === 0) {
            items.items.push({
                'title': 'Haltestelle nicht gefunden 🤔',
                'subtitle': 'Vielleicht ein Tippfehler?'
            })
            return items
        }

        items.items = data.map(function (con) {
            var timeText
            if (con.arrivalTimeRelative === 0) {
                timeText = ' jetzt'
            } else if (con.arrivalTimeRelative === 1) {
                timeText = ' in 1 Minute'
            } else {
                timeText = ' in ' + con.arrivalTimeRelative + ' Minuten'
            }
            return {
                'title': con.line + ' ' + con.direction + timeText,
                'subtitle': moment().add(con.arrivalTimeRelative, 'm').format('dddd, HH:mm [Uhr]'),
                'icon': {
                    'path': 'transport_icons/' + con.mode.name + '.png'
                }
            }
        })
        return items
    })
    .catch(function (err) {
        var items = {'items': [{
            'title': 'Unerwarteter Fehler 😲',
            'subtitle': err.message,
        }]}
        return items
    })
}
