## alfred dvb

Monitor public transport in the VVO/DVB network right from within Alfred. Powered by [dvbgo](https://github.com/kiliankoe/dvbgo).

![screenshot](https://cloud.githubusercontent.com/assets/2625584/17901215/b9d2f916-6962-11e6-8c34-665176f36a82.png)

Download [here](https://github.com/kiliankoe/alfred_dvb/releases/latest) (click on `DVB.v*.alfredworkflow`).

### Usage

- `dvb helmholtz`

  Gives you all upcoming connections from Helmholtzstraße.

- `dvb prager in 10`

  Gives you all upcoming connections from Prager Straße in 10 minutes. Optional text after the amount of minutes is ignored, so you could also enter `dvb prager in 10 minutes` if you prefer.

- `dvb albertplatz [3]`

  Gives you all upcoming connections from Albertplatz that are serviced by the line 3. 

You can of course mix and match time offsets with `... in x` and filters with `[x]` in one query. Whichever comes first should also be completely up to you, just be sure to put time offsets after the name of the stop, e.g. `dvb pirnaischer [62] in 10` or `dvb hbf in 60 [s3]` are valid.

Hitting ⏎ on a connection will schedule a notification to be sent 10 minutes prior to departure. This obviously only works for departures that will depart more than 10 minutes from the current time. 

### Problems?

Please [report an issue](https://github.com/kiliankoe/alfred_dvb/issues/new) if something isn't working as expected or you have a question/feature request.

### Credits

Bus icon by [icons8.com](https://icons8.com).
