package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	gorez "github.com/JackStillwell/GoRez/pkg"
)

func main() {
	var authPath, dataDirPath, matchIds string
	var getGods, getItems bool
	flag.StringVar(&authPath, "auth", "", "The file path to the hirez dev auth file")
	flag.StringVar(&dataDirPath, "datadir", "",
		"The file path to the directory containing SMITE data")
	flag.StringVar(&matchIds, "matchids", "", "A CSV list of matchids to retrieve")
	flag.BoolVar(&getGods, "gods", false, "Fetch all gods")
	flag.BoolVar(&getItems, "items", false, "Fetch all items")

	flag.Parse()

	if authPath == "" {
		log.Fatal("'auth' argument is required")
		flag.PrintDefaults()
	}
	if dataDirPath == "" {
		log.Fatal("'datadir' argument is required")
		flag.PrintDefaults()
	}

	log.Println("instantiating gorez")
	g, err := gorez.NewGorez(authPath)
	if err != nil {
		log.Fatal("failed to instantiate gorez: ", err)
	}
	defer g.Shutdown()
	log.Println("gorez instantiated")

	log.Println("initing gorez")
	err = g.Init()
	if err != nil {
		log.Fatal("failed to init gorez: ", err)
	}
	log.Println("gorez inited")

	if getGods {
		gods, err := g.GetGods()
		if err != nil {
			log.Println("error fetching gods: ", err)
		} else {
			jBytes, err := json.Marshal(gods)
			if err != nil {
				log.Println("error marshaling gods", err)
			}
			godsPath := path.Join(dataDirPath, "gods.json")
			f, err := os.Create(godsPath)
			if err != nil {
				log.Println("error opening file to write gods", err)
			}
			f.Close()
			f.Write(jBytes)
		}
	}

	if getItems {
		items, err := g.GetItems()
		if err != nil {
			log.Println("error fetching items", err)
		} else {
			jBytes, err := json.Marshal(items)
			if err != nil {
				log.Println("error marshaling items", err)
			}
			itemsPath := path.Join(dataDirPath, "items.json")
			f, err := os.Create(itemsPath)
			if err != nil {
				log.Println("error opening file to write items", err)
			}
			f.Close()
			f.Write(jBytes)
		}
	}

	if matchIds != "" {
		idStrings := strings.Split(matchIds, ",")
		idInts := make([]int, 0, len(idStrings))
		for _, s := range idStrings {
			intId, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Println("error parsing matchid", s, "to int")
				continue
			}
			idInts = append(idInts, int(intId))
		}

		matchDetails, errs := g.GetMatchDetailsBatch(idInts...)
		log.Println("errors fetching items", errs)
		jBytes, err := json.Marshal(matchDetails)
		if err != nil {
			log.Println("error marshaling matchdetails", err)
		}
		matchDetailsPath := path.Join(dataDirPath, "matchdetails.json")
		f, err := os.Create(matchDetailsPath)
		if err != nil {
			log.Println("error opening file to write matchdetails", err)
		}
		f.Close()
		f.Write(jBytes)
	}

}
