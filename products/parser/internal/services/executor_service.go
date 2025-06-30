package services

import (
	"context"
	"fmt"
	"log"
	"scrappers/internal/domain"
	"scrappers/internal/helpers"
	"time"
)

type ExecutorService struct {
	processor executorHelper
	config    domain.IConfigService
}

type result struct {
	ProductInfos []domain.ProductInfo
	ShopName     string
	Error        error
}

type executorHelper struct {
	config        domain.IConfigService
	matcher       domain.IMatcherService
	writerFactory domain.IWriterFactory
	lenta         *helpers.LentaRequester
	perekrestok   *helpers.PerekrestokRequester
	dixy          *helpers.DixyRequester
	magnit        *helpers.MagnitRequester
}

func NewExecutorService(config domain.IConfigService, matcher domain.IMatcherService, writerFactory domain.IWriterFactory, lenta *helpers.LentaRequester, perekrestok *helpers.PerekrestokRequester, dixy *helpers.DixyRequester, magnit *helpers.MagnitRequester) *ExecutorService {
	return &ExecutorService{processor: executorHelper{config: config, matcher: matcher, writerFactory: writerFactory, lenta: lenta, perekrestok: perekrestok, dixy: dixy, magnit: magnit}, config: config}
}

func (service executorHelper) Process() error {
	config, err := service.config.Get()
	if err != nil {
		return fmt.Errorf("can't get config: %+v", err)
	}
	writer, err := service.writerFactory.Get()
	if err != nil {
		return fmt.Errorf("can't get writer service from factory: %+v", err)
	}

	newRequestDelay := config.DelayInfo.NewRequestDelay

	for _, category := range config.Categories {
		ch := make(chan result)

		total := 0
		for _, lenta := range category.Lenta {
			go service.processLenta(ch, category.Type, lenta, newRequestDelay)
			total += 1
		}
		for _, perekrestok := range category.Perekrestok {
			go service.processPerekrestok(ch, category.Type, perekrestok, newRequestDelay)
			total += 1
		}
		for _, dixy := range category.Dixy {
			go service.processDixy(ch, category.Type, dixy, newRequestDelay)
			total += 1
		}
		for _, magnit := range category.Magnit {
			go service.processMagnit(ch, category.Type, magnit, newRequestDelay)
			total += 1
		}

		products := []domain.ProductInfo{}

		log.Printf("Processor: start collecting products for category=%s", category.Type)

		for i := 0; i < total; i += 1 {
			data := <-ch
			if data.Error != nil {
				log.Printf("Processor: can't get %s items for category=%s, error=%s", data.ShopName, category.Type, data.Error)
			} else {
				products = append(products, data.ProductInfos...)
			}
		}

		match := service.matcher.Match(products, category.Type)

		if err := writer.Write(match); err != nil {
			log.Printf("can't write data: %s", err)
		}
	}

	log.Println("scrapped all data")

	return writer.Close()
}

func (service executorHelper) Ping() {
	lentaErr := service.lenta.Ping()
	perekrestokErr := service.perekrestok.Ping()
	dixyErr := service.dixy.Ping()
	magnitErr := service.magnit.Ping()

	if lentaErr != nil {
		log.Printf("error ping lenta: %s", lentaErr)
	}
	if perekrestokErr != nil {
		log.Printf("error ping perekrestok: %s", perekrestokErr)
	}
	if dixyErr != nil {
		log.Printf("error ping dixy: %s", dixyErr)
	}
	if magnitErr != nil {
		log.Printf("error ping magnit: %s", magnitErr)
	}
}

func (service executorHelper) processLenta(ch chan<- result, categoryType string, category domain.LentaCategory, newRequestDelay int) {
	infos, err := service.lenta.Process(category.Id, category.Multicheckboxes, newRequestDelay)
	log.Printf("Lenta: scrapped %s, len=%d", categoryType, len(infos))
	if err != nil {
		log.Printf("Lenta: category=%s, err=%+v", categoryType, err)
		ch <- result{Error: err}
	} else {
		log.Printf("Lenta: category=%s, err=nil", categoryType)
		ch <- result{ProductInfos: infos, ShopName: "Лента"}
	}
	log.Printf("Lenta: exit from category=%s", categoryType)
}

func (service executorHelper) processPerekrestok(ch chan<- result, categoryType string, category domain.PerekrestokCategory, newRequestDelay int) {
	infos, err := service.perekrestok.Process(category.Id, category.Features, newRequestDelay)
	log.Printf("Perekrestok: scrapped %s, len=%d", categoryType, len(infos))
	if err != nil {
		log.Printf("Perekrestok: category=%s, err=%+v", categoryType, err)
		ch <- result{Error: err}
	} else {
		log.Printf("Perekrestok: category=%s, err=nil", categoryType)
		ch <- result{ProductInfos: infos, ShopName: "Перекрёсток"}
	}
	log.Printf("Perekrestok: exit from category=%s", categoryType)
}

func (service executorHelper) processDixy(ch chan<- result, categoryType string, category domain.DixyCategory, newRequestDelay int) {
	infos, err := service.dixy.Process(category.Id, category.Filters, newRequestDelay)
	log.Printf("Dixy: scrapped %s, len=%d", categoryType, len(infos))
	if err != nil {
		log.Printf("Dixy: category=%s, err=%+v", categoryType, err)
		ch <- result{Error: err}
	} else {
		log.Printf("Dixy: category=%s, err=nil", categoryType)
		ch <- result{ProductInfos: infos, ShopName: "Дикси"}
	}
	log.Printf("Dixy: exit from category=%s", categoryType)
}

func (service executorHelper) processMagnit(ch chan<- result, categoryType string, category domain.MagnitCategory, newRequestDelay int) {
	infos, err := service.magnit.Process(category.Id, category.Filters, newRequestDelay)
	log.Printf("Magnit: scrapped %s, len=%d", categoryType, len(infos))
	if err != nil {
		log.Printf("Magnit: category=%s, err=%+v", categoryType, err)
		ch <- result{Error: err}
	} else {
		log.Printf("Magnit: category=%s, err=nil", categoryType)
		ch <- result{ProductInfos: infos, ShopName: "Магнит"}
	}
	log.Printf("Magnit: exit from category=%s", categoryType)
}

func (service *ExecutorService) Start() {
	log.Printf("start scrapping: %d", time.Now().UnixMilli())
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Hour)

		if err := service.processor.Process(); err != nil {
			log.Printf("error scrapping products: %+v", err)
		}
		if err := service.ping(ctx); err != nil {
			log.Printf("error ping sites: %+v", err)
		}

		cancel()
	}
}

func (service *ExecutorService) ping(ctx context.Context) error {
	log.Printf("start pinging: %d", time.Now().UnixMilli())

	config, err := service.config.Get()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(config.DelayInfo.PingDelay) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			log.Printf("pinging: %d", time.Now().UnixMilli())
			service.processor.Ping()
		}
	}
}
