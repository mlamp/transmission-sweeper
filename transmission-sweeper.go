package main

import (
	"flag"
	"strconv"
	"time"
    
	log "github.com/Sirupsen/logrus"
	"github.com/tubbebubbe/transmission"
	"github.com/vharitonsky/iniflags"
)

var (
	flagProtocol = flag.String("flagProtocol", "http", "Transmission server protocol")
	flagHost     = flag.String("flagHost", "localhost", "Transmission server host")
	flagPort     = flag.Int("flagPort", 9091, "Transmission server port")

	flagUsername = flag.String("flagUsername", "", "Transmission RPC username")
	flagPassword = flag.String("flagPassword", "", "Transmission RPC password")

	flagTorrentOlderThanDays = flag.Int("flagTorrentOlderThanDays", 7, "Torrent is older than days")
    flagRatioLowerThan = flag.Float64("flagRatioLowerThan", -1, "Torrent have lower ratio than")
    flagErrorFilter = flag.Bool("flagErrorFilter", false, "Filter torrents with error / non-error, if not set dont filter at all")
    flagHasError = flag.Bool("flagHasError", false, "Torrent has error (status > 0)")
    flagDryRun = flag.Bool("flagDryRun", false, "Script runs in simulation mode, no actual deletion")
)

func main() {
	iniflags.Parse()
	client := transmission.New(*flagProtocol+"://"+*flagHost+":"+strconv.Itoa(*flagPort), *flagUsername, *flagPassword)

	torrents, err := client.GetTorrents()
	if err != nil {
		log.Panic(err)
	}

	torrents = applyTimeFilter(torrents, *flagTorrentOlderThanDays)

    if *flagErrorFilter {
        torrents = applyErrorFilter(torrents, *flagHasError)
    }
    
    if *flagRatioLowerThan >= 0.0 {
        torrents = applyRatioFilter(torrents, *flagRatioLowerThan)
    }

	for _, torrent := range torrents {
        log.WithFields(log.Fields{
            "ID": torrent.ID,
            "Name": torrent.Name,
            "Status": torrent.Status,
            "LeftUntilDone": torrent.LeftUntilDone,
            "AddedDate": torrent.AddedDate,
            "Eta": torrent.Eta,
            "UploadRatio": torrent.UploadRatio,
            "RateDownload": torrent.RateDownload,
            "RateUpload": torrent.RateUpload,
            "DownloadDir": torrent.DownloadDir,
            "IsFinished": torrent.IsFinished,
            "PercentDone": torrent.PercentDone,
            "SeedRatioMode": torrent.SeedRatioMode,
            "Error": torrent.Error,
            "ErrorString": torrent.ErrorString,
        }).Info("Torrent going for deletion")
        
        if !*flagDryRun {
            transmissionRemoveCmd, _ := transmission.NewDelCmd(torrent.ID, true)
            _, err := client.ExecuteCommand(transmissionRemoveCmd)
            
            if err != nil {
                log.WithFields(log.Fields{
                    "ID": torrent.ID,
                    "Name": torrent.Name,
                    "Status": torrent.Status,
                    "LeftUntilDone": torrent.LeftUntilDone,
                    "AddedDate": torrent.AddedDate,
                    "Eta": torrent.Eta,
                    "UploadRatio": torrent.UploadRatio,
                    "RateDownload": torrent.RateDownload,
                    "RateUpload": torrent.RateUpload,
                    "DownloadDir": torrent.DownloadDir,
                    "IsFinished": torrent.IsFinished,
                    "PercentDone": torrent.PercentDone,
                    "SeedRatioMode": torrent.SeedRatioMode,
                    "Error": torrent.Error,
                    "ErrorString": torrent.ErrorString,
                }).Error("Could'nt delete torrent")
            }
        }
	}
}

func applyTimeFilter(t transmission.Torrents, olderThanDays int) transmission.Torrents {
	var filtered transmission.Torrents
	timeDuration := time.Duration(int64(*flagTorrentOlderThanDays)) * 24 * time.Hour
	for _, torrent := range t {
		addedDate := time.Unix(int64(torrent.AddedDate), 0)
		if time.Now().Add(-timeDuration).After(addedDate) {
			filtered = append(filtered, torrent)
		}
	}
	return filtered
}


func applyRatioFilter(t transmission.Torrents, ratioLowerThan float64) transmission.Torrents {
	var filtered transmission.Torrents
	for _, torrent := range t {
        if torrent.UploadRatio < ratioLowerThan {
            filtered = append(filtered, torrent)
        }
	}
	return filtered
}


func applyErrorFilter(t transmission.Torrents, hasError bool) transmission.Torrents {
	var filtered transmission.Torrents
	for _, torrent := range t {
		if hasError && torrent.Error > 0 {
			filtered = append(filtered, torrent)
		} else if !hasError && torrent.Error == 0 {
            filtered = append(filtered, torrent)
        }
	}
	return filtered
}

