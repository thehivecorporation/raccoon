package book

import (
	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/instructions"
	"encoding/json"
	"io/ioutil"
)

type BOOK struct {
	RUN []instructions.RUN
	ADD []instructions.ADD
}
func (z *BOOK) ReadZbook(zombie_book string) BOOK {
	return read(*z,zombie_book)
}

func read(z BOOK, zombie_book string) BOOK{

	log.Info("Reading file -----------------> " + zombie_book)
	dat, err := ioutil.ReadFile(zombie_book)
	if err != nil {
		log.Error(err)
	}
	json.Unmarshal([]byte(string(dat)), &z)
	return z

}
