package api

import (
	"encoding/json"
	"fmt"
	"kibble/models"
	"sort"
	"strings"

	"github.com/shopspring/decimal"
)

// LoadAllPrices will load all prices
func LoadAllPrices(cfg *models.Config, site *models.Site, itemIndex models.ItemIndex) error {

	if cfg.DefaultPricingCountryCode == "" {
		log.Info("skipping pricing, not country code specified")
		return nil
	}

	slugs := make([]string, 0)

	for k := range itemIndex["film"] {
		slugs = append(slugs, k)
	}

	for k := range itemIndex["bundle"] {
		slugs = append(slugs, k)
	}

	for k := range itemIndex["tv-season"] {
		slugs = append(slugs, k)
	}

	for k := range itemIndex["plan"] {
		slugs = append(slugs, k)
	}

	sort.Strings(slugs)

	const batchSize = 300
	var total = 0
	var s []string
	for len(slugs) > 0 {

		if len(slugs) > batchSize {
			s = slugs[:batchSize]
			slugs = slugs[batchSize:]
		} else {
			s = slugs[:]
			slugs = nil
		}

		count, err := loadPrices(cfg, site, s, itemIndex)
		if err != nil {
			return err
		}
		total += count
	}

	log.Infof("prices: loaded %d for %s", total, cfg.DefaultPricingCountryCode)

	return nil
}

func loadPrices(cfg *models.Config, site *models.Site, slugs []string, itemIndex models.ItemIndex) (int, error) {

	ids := strings.Join(slugs, ",")
	path := fmt.Sprintf("%s/services/pricing/v2/prices/show_multiple?items=%s&location=%s", cfg.SiteURL, ids, cfg.DefaultPricingCountryCode)

	data, err := Get(cfg, path)
	if err != nil {
		log.Infof("pricing failed to load %s", err)
		return 0, err
	}

	var details prices
	err = json.Unmarshal([]byte(data), &details)
	if err != nil {
		log.Error("price.error: %s", err)
		log.Debug("invalid data %s", string(data))
		return 0, err
	}

	return processPrices(details, site, itemIndex)
}

func processPrices(details prices, site *models.Site, itemIndex models.ItemIndex) (int, error) {

	count := 0
	for _, p := range details.Prices {
		found := itemIndex.Get(p.Item)

		pricingInfo := p.getPrices()

		switch found.ItemType {
		case "film":
			if f, ok := site.Films.FindFilmBySlug(p.Item); ok {
				count++
				f.PriceInfo = pricingInfo
				// replace the itemIndex
				itemIndex.Replace(p.Item, f.GetGenericItem())
			}
		case "tvseason":
			if f, ok := site.TVSeasons.FindTVSeasonBySlug(p.Item); ok {
				count++
				f.PriceInfo = pricingInfo
				// replace the itemIndex
				itemIndex.Replace(p.Item, f.GetGenericItem())
			}

		case "bundle":
			if f, err := site.Bundles.FindBundleBySlug(p.Item); err == nil {
				count++
				f.PriceInfo = pricingInfo
				// replace the itemIndex
				itemIndex.Replace(p.Item, f.GetGenericItem())
			}
		case "plan":
			if f, err := site.Plans.FindPlanBySlug(p.Item); err == nil {
				count++
				f.PriceInfo = pricingInfo
				// replace the itemIndex
				itemIndex.Replace(p.Item, f.GetGenericItem())
			}
		}
	}

	// loop over the retrieved list of films that belong to plans
	// assign a slice of maps to each film e.g. [{"plan/123" : prices }, ...]
	for _, filmPlans := range details.Plans {
		if film, ok := site.Films.FindFilmBySlug(filmPlans.Item); ok {

			filmPlanPrices := make(map[string]models.PriceCollection, 0)

			for _, filmPlan := range filmPlans.Plans {
				if plan, err := site.Plans.FindPlanBySlug(filmPlan); err == nil {
					filmPlanPrices[filmPlan] = plan.PriceInfo.Prices
				}
			}
			film.PriceInfo.PlanPrices = filmPlanPrices
		}
	}

	return count, nil
}

type prices struct {
	Prices []pricesV2 `json:"prices"`
	Plans  []struct {
		Item  string   `json:"item"`
		Plans []string `json:"plans"`
	} `json:"plans"`
}

type pricesV2 struct {
	Item     string          `json:"item"`
	Currency string          `json:"currency"`
	Rent     *qualityPriceV2 `json:"rent"`
	Buy      *qualityPriceV2 `json:"buy"`
}

type qualityPriceV2 struct {
	Hd       *string `json:"hd"`
	HdString *string `json:"hd_string"`
	Sd       *string `json:"sd"`
	SdString *string `json:"sd_string"`
}

func (p pricesV2) getPrices() models.PriceInfo {
	pp := make([]models.Price, 0)

	if p.Buy != nil {
		pp = append(pp, p.Buy.getPrice(models.Buy)...)
	}
	if p.Rent != nil {
		pp = append(pp, p.Rent.getPrice(models.Rent)...)
	}

	return models.PriceInfo{Currency: p.Currency, Prices: pp}
}

func (qp qualityPriceV2) getPrice(ownership string) []models.Price {
	pp := make([]models.Price, 0)

	if qp.Hd != nil && qp.HdString != nil {
		price, err := decimal.NewFromString(*qp.Hd)
		if err == nil {
			pp = append(pp, models.Price{Ownership: ownership, Quality: models.HD, Price: price, PriceString: *qp.HdString})
		}
	}

	if qp.Sd != nil && qp.SdString != nil {
		price, err := decimal.NewFromString(*qp.Sd)
		if err == nil {
			pp = append(pp, models.Price{Ownership: ownership, Quality: models.SD, Price: price, PriceString: *qp.SdString})
		}
	}

	return pp
}
