# Retentor
*DigitalOcean Spaces Backup Retention Policy Agent*

Retentor keeps track of all objects in the specified DigitalOcean Space and manages their retention. It expects backups to be contained within folders specifying date and time. Eg: /root/(yyyymmdd)/(HHMM)/


## Policy Types
Retentor has a few different policies types:

- `hourly`
- `sixhourly`
- `daily`
- `weekly`
- `monthly`

These policies are enforced within the specified timeframe.

**Important**: If, for example, a monthly policy is specified, but the timeframe does not include the start of a month, all backups within that timeframe are deleted.

### hourly
Keeps the oldest backup for every hour
### sixhourly
Keeps the oldest backup for every six hours
### daily
Keeps the oldest backup for every day
### weekly
Keeps the oldest backup of the last day of a every week in the timeframe or/and the oldest backup of the last day of the oldest week in the timefame
### monthly
Keeps the oldest backup of the last day of a every month in the timeframe or/and the oldest backup of the last day of the oldest month in the timeframe

## Defining Policies
The `policies.go` file specifies the retention policies