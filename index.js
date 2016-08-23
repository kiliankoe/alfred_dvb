#! /usr/local/bin/node

var dvb = require('dvbjs')
var moment = require('moment')
moment.locale('de')

var args = process.argv.slice(2)

var timeOffset
args[1] ? timeOffset = args[1] : timeOffset = 0

var numResults = 6
dvb.monitor(args[0], timeOffset, numResults)
    .then(function(data) {
        var items = {'items': []}

        if (data.length === 0) {
            items.items.push({
                'title': 'Haltestelle nicht gefunden 🤔'
            })
            console.log(JSON.stringify(items))
            return
        }

        data.forEach(function(con) {
            var timeText
            if (con.arrivalTimeRelative === 0) {
                timeText = ' jetzt'
            } else if (con.arrivalTimeRelative === 1) {
                timeText = ' in 1 Minute'
            } else {
                timeText = ' in ' + con.arrivalTimeRelative + ' Minuten'
            }
            items.items.push({
                'title': con.line + ' ' + con.direction + timeText,
                'subtitle': moment().add(con.arrivalTimeRelative, 'm').format('dddd, HH:mm [Uhr]'),
                'icon': {
                    'path': 'transport_icons/' + con.mode.name + '.png'
                }
            })
        })

        console.log(JSON.stringify(items))
    })
    .catch(function (err) {
        var items = {'items': [{
            'title': 'Unerwarteter Fehler 😲',
            'subtitle': err.message,
        }]}
        console.log(JSON.stringify(items))
        throw err
    })
