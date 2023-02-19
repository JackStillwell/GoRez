package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"

	gorez "github.com/JackStillwell/GoRez/pkg"
)

func main() {
	var authPath, dataDirPath string
	flag.StringVar(&authPath, "auth", "", "The file path to the hirez dev auth file")
	flag.StringVar(&dataDirPath, "datadir", "",
		"The file path to the directory containing SMITE data")

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
	log.Println("gorez instantiated")

	log.Println("initing gorez")
	err = g.Init()
	if err != nil {
		log.Fatal("failed to init gorez: ", err)
	}
	log.Println("gorez inited")

	gods, err := g.GetGods()
	if err != nil {
		log.Println("error fetching gods: ", err)
	} else {
		jString, err := json.Marshal(gods)
		if err != nil {
			log.Println("error marshaling gods: ", err)
		}
		godsPath := path.Join(dataDirPath, "gods.json")
		err = os.WriteFile(godsPath, []byte(jString),
			os.FileMode(os.O_CREATE|os.O_TRUNC|os.O_RDWR))
		if err != nil {
			log.Println("error writing gods: ", err)
		}
	}

	items, err := g.GetItems()
	if err != nil {
		log.Println("error fetching items", err)
	} else {
		jString, err := json.Marshal(items)
		if err != nil {
			log.Println("error marshaling items", err)
		}
		itemsPath := path.Join(dataDirPath, "items.json")
		err = os.WriteFile(itemsPath, []byte(jString),
			os.FileMode(os.O_CREATE|os.O_TRUNC|os.O_RDWR))
		if err != nil {
			log.Println("error writing items", err)
		}
	}
}
