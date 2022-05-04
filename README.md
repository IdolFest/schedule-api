# schedule-api

Schedule API parses a specially formatted Google Sheet and produces usable JSON. It is compatible
with the [PonyFest grid renderer](https://github.com/PonyFest/website/blob/master/static/scripts/schedule.js)
as well as all the other PonyFest tools that consume schedule inforamtion.

## Usage

```
  --allowed-origins string
    	The set of Origins that should be returned for requests. "*" is a good choice, but not default.
  --bind string
    	The host:port to bind to. (default "0.0.0.0:8080")
  --cache-timeout int
    	The timeout in seconds when a new copy of the schedule should be fetched. This applies also when the schedule cannot be fetched. (default 300)
  --row-minutes int
    	The number of minutes represented by one row. (default 30)
  --sheet-url string
    	The URL of the published Schedule Spreadsheet. Expected response is in CSV format.
  --timezone string
    	The timezome to assume for the spreadsheet (default "America/New_York")
```

## Google Sheet format

The tool expects to consume the CSV format of a Google Sheet. To create this, go to File -> Share ->
Publish To Web, then select the relevant sheet and set the format to "Comma-separated values (.csv)".
Provide the URL Sheets provides you with to `--sheet-url` - it should end with `&output=csv`.

The google sheet is expected to have two header rows: one with the room name, and a second that
contains repeating column headers for each room. Currently, the per-room columns must be:

* `ID` - a unique per-event ID exposed in the API. This is also the field used to determine how long
  each event is. It is opaque, and can be any unique string.
* `Title` - The title of the event
* `Panelists` - The panelists for the event
* `Public Description` - A public description of the event
* `Is Zoom` - Whether the event is taking place via Zoom.

Additionally, one column on the left must contain the event times. in YYYY-MM-DD HH:MM format. Each
row must represent a consistent length of time.

As a sample,
[this was the schedule for PonyFest Online! 4.0](https://docs.google.com/spreadsheets/d/1UHjn4SEqcmZjfXd1VoHLzP1UG0HFO5BKMpFFxLsjzyY/edit?usp=sharing), which rendered to [this JSON](https://ponyfest.horse/4.0/schedule.json).
