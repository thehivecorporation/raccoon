package book

import (
	log "github.com/Sirupsen/logrus"
	"github.com/thehivecorporation/raccoon/instructions"

)

type BOOK struct {
	RUN []instructions.RUN
	ADD []instructions.ADD
}

func (b *BOOK) readBook(zombie_book string) {
	read(b,zombie_book)

}

func read(b BOOK, zombie_book string) {
	log.Info("Reading file ---------------------------------> " + zombie_book)


}






