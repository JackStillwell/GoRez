package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	gorez "github.com/JackStillwell/GoRez/pkg"
	models "github.com/JackStillwell/GoRez/pkg/models"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		fmt.Println("failed to init logger: ", err)
	}

	var authPath, dataDirPath, matchIds string
	var numDays int
	var getGods, getItems bool
	flag.StringVar(&authPath, "auth", "", "The file path to the hirez dev auth file")
	flag.StringVar(&dataDirPath, "datadir", "",
		"The file path to the directory containing SMITE data")
	flag.StringVar(&matchIds, "matchids", "", "A CSV list of matchids to retrieve")
	flag.IntVar(&numDays, "numdays", 0,
		"The number of days into the past to retrieve RankedConquest matches for")
	flag.BoolVar(&getGods, "gods", false, "Fetch all gods")
	flag.BoolVar(&getItems, "items", false, "Fetch all items")

	flag.Parse()

	if authPath == "" {
		log.Error("'auth' argument is required")
		flag.PrintDefaults()
		return
	}
	if dataDirPath == "" {
		log.Error("'datadir' argument is required")
		flag.PrintDefaults()
		return
	}

	log.Debug("instantiating gorez")
	g, err := gorez.NewGorez(authPath, nil, 10, log)
	if err != nil {
		log.Error("failed to instantiate gorez", zap.Error(err))
		return
	}
	defer g.Shutdown()
	log.Debug("gorez instantiated")

	log.Debug("initing gorez")
	err = g.Init()
	if err != nil {
		log.Error("failed to init gorez", zap.Error(err))
		return
	}
	log.Debug("gorez inited")

	if getGods {
		gods, err := g.GetGods()
		if err != nil {
			log.Error("error fetching gods", zap.Error(err))
		} else {
			jBytes, err := json.Marshal(gods)
			if err != nil {
				log.Error("error marshaling gods", zap.Error(err))
			} else {
				godsPath := path.Join(dataDirPath, "gods.json")
				f, err := os.Create(godsPath)
				if err != nil {
					log.Error("error opening file to write gods", zap.Error(err))
					return
				} else {
					f.Close()
					f.Write(jBytes)
				}
			}
		}
	}

	if getItems {
		items, err := g.GetItems()
		if err != nil {
			log.Error("error fetching items", zap.Error(err))
		} else {
			jBytes, err := json.Marshal(items)
			if err != nil {
				log.Error("error marshaling items", zap.Error(err))
			} else {
				itemsPath := path.Join(dataDirPath, "items.json")
				f, err := os.Create(itemsPath)
				if err != nil {
					log.Error("error opening file to write items", zap.Error(err))
				} else {
					f.Close()
					f.Write(jBytes)
				}
			}
		}
	}

	if matchIds != "" {
		idStrings := strings.Split(matchIds, ",")
		matchDetails, errs := g.GetMatchDetailsBatch(idStrings...)
		log.Error("failed fetching matchdetailsbatch", zap.Errors("errors", errs))
		jBytes, err := json.Marshal(matchDetails)
		if err != nil {
			log.Error("failed marshaling matchdetails", zap.Error(err))
			return
		}
		matchDetailsPath := path.Join(dataDirPath, "matchdetails.json")
		f, err := os.Create(matchDetailsPath)
		if err != nil {
			log.Error("failed opening file to write matchdetails", zap.Error(err))
			return
		}
		defer f.Close()
		nBytes, err := f.Write(jBytes)
		if err != nil {
			log.Error("failed writing matchdetails file:", zap.Error(err))
			return
		}
		if nBytes == 0 {
			log.Error("no bytes written to matchdetails file")
		}
	}

	if numDays != 0 {
		queueIDs := []models.QueueID{models.RankedConquest}
		dateStrings := make([]string, 0, numDays)
		currDate := time.Now()
		for i := 0; i < numDays; i++ {
			year := currDate.Year()
			month := currDate.Month()
			day := currDate.Day()

			dateStrings = append(dateStrings, fmt.Sprintf("%d%02d%02d/0,00", year, month, day))
			currDate = currDate.Add(-24 * time.Hour)
		}

		matchIdsPath := path.Join(dataDirPath, "matchids.json")
		f, err := os.Create(matchIdsPath)
		if err != nil {
			log.Error("error opening file to write matchids", zap.Error(err))
			return
		}
		defer f.Close()

		matchIds, errs := g.GetMatchIDsByQueue(dateStrings, queueIDs)
		log.Error("failed fetching matchidsbyqueue", zap.Errors("errors", errs))
		jBytes, err := json.Marshal(matchIds)
		if err != nil {
			log.Error("failed marshaling matchids", zap.Error(err))
			return
		}

		nBytes, err := f.Write(jBytes)
		if err != nil {
			log.Error("failed writing matchids file", zap.Error(err))
			return
		}
		if nBytes == 0 {
			log.Error("no bytes written to matchids file")
			return
		}

		toRetrieve := []string{}
		for _, matchIdList := range matchIds {
			if matchIdList == nil {
				continue
			}
			for _, matchId := range *matchIdList {
				if matchId.Match != nil {
					toRetrieve = append(toRetrieve, *matchId.Match)
				}
			}
		}

		log.Info("retrieving matchids", zap.Strings("matchids", toRetrieve))
		bytesList, errs := g.GetMatchDetailsBatch(toRetrieve...)
		log.Error("failures fetching matchdetailsbatch", zap.Errors("errors", errs))

		dateString := time.Now().UTC().Format("2006-Jan-02")
		for i, bytes := range bytesList {
			matchIdsPath := path.Join(dataDirPath, fmt.Sprintf("matchdetails-%s_%d.json",
				dateString, i))
			f, err := os.Create(matchIdsPath)
			if err != nil {
				log.Error("failed opening file to write matchdetails", zap.Int("fileNumber", i),
					zap.Error(err))
				return
			}
			defer f.Close()
			nBytes, err := f.Write(bytes)
			if err != nil {
				log.Error("failed writing matchdetails", zap.Int("fileNumber", i),
					zap.Error(err))
				return
			}
			if nBytes == 0 {
				log.Error("no bytes written to matchdetails", zap.Int("fileNumber", i))
			}
		}
	}
}
