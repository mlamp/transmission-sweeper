# Transmission Sweeper - tool to clean up old torrents

## Installation

```bash
go get github.com/mlamp/transmission-sweeper
```

## Usage

```bash
bin/transmission-sweeper
```

```text
  -flagDryRun
        Script runs in simulation mode, no actual deletion
  -flagErrorFilter
        Filter torrents with error / non-error, otherwise all
  -flagHasError
        Torrent has error
  -flagHost string
        Transmission server host (default "localhost")
  -flagPassword string
        Transmission RPC password
  -flagPort int
        Transmission server port (default 9091)
  -flagProtocol string
        Transmission server protocol (default "http")
  -flagRatioLowerThan float
        Torrent have lower ratio than (default -1)
  -flagTorrentOlderThanDays int
        Torrent is older than days (default 7)
  -flagUsername string
        Transmission RPC username
```

## Examples

Simulate deletion of torrents older than 0 days and having ratio lower than 0.1

```bash
transmission-sweeper -flagHost=localhost -flagPort=9091 -flagUsername=transmission -flagPassword=secret123 -flagTorrentOlderThanDays=0 -flagRatioLowerThan=0.1 -flagDryRun=false
```

Actually delete torrents (missing -flagDryRun) older than 30 days and having ratio lower than 3

```bash
transmission-sweeper -flagHost=localhost -flagPort=9091 -flagUsername=transmission -flagPassword=secret123 -flagTorrentOlderThanDays=30 -flagRatioLowerThan=3
```

Actually delete torrents which are more than 10 days old and have error status

```bash
transmission-sweeper -flagHost=localhost -flagPort=9091 -flagUsername=transmission -flagPassword=secret123 -flagTorrentOlderThanDays=30 -flagErrorFilter=true -flagHasError=true
```
