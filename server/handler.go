package server

import (
	"errors"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/jeffreyyong/xe/calculator"
	"github.com/jeffreyyong/xe/client"
	"github.com/jeffreyyong/xe/date"
	"github.com/jeffreyyong/xe/model"
)

const (
	ParamCurrency = "currency"

	// Number of days before the current date for historical rates
	DaysForRates = 7
)

var noRouteFoundFunc = func(c *gin.Context) {
	c.JSON(http.StatusNotFound, &model.ConvertResp{Error: model.ErrRouteNotFound})
}

// Handler that has forex client
// and calculator engine
type Handler struct {
	fx client.Forex
	ce calculator.Engine
}

// NewHandler initialises a Handler
func NewHandler(forex client.Forex, calculator calculator.Engine) *Handler {
	return &Handler{
		fx: forex,
		ce: calculator,
	}
}

// SetupAPIHandler sets up a GIN router
// with /convert GET endpoint
func SetupAPIHandler(h *Handler) *gin.Engine {
	r := gin.Default()
	r.NoRoute(noRouteFoundFunc)
	r.GET(model.ConvertEndpoint, h.Convert)
	return r
}

// Convert is the handler func for /convert endpoint
func (h *Handler) Convert(ctx *gin.Context) {
	httpStatus, forexResp, err := h.convert(ctx)
	if err != nil {
		log.Print(err)
	}
	ctx.JSON(httpStatus, forexResp)
}

func (h *Handler) convert(ctx *gin.Context) (int, *model.ConvertResp, error) {
	currency := ctx.Query(ParamCurrency)
	if currency == "" {
		return http.StatusBadRequest, &model.ConvertResp{Error: model.ErrDecodeParams}, nil
	}

	// get latest rate
	latestRate, err := h.fx.GetLatestRate(currency)
	if err != nil {
		return http.StatusInternalServerError, &model.ConvertResp{Error: model.ErrConvert}, err
	}

	// extract the rate
	targetRate, err := extractTargetRate(latestRate)
	if err != nil {
		return http.StatusInternalServerError, &model.ConvertResp{Error: model.ErrConvert}, err
	}

	// compute the recommendation
	recommendation, err := h.computeRecommendation(currency)
	if err != nil {
		return http.StatusInternalServerError, &model.ConvertResp{Error: model.ErrConvert}, err
	}

	convertResp := &model.ConvertResp{
		From:           currency,
		To:             calculator.EUR,
		Rate:           targetRate,
		Recommendation: string(recommendation),
	}
	return http.StatusOK, convertResp, nil
}

// computeRecommendation
// 1. generates a start and end date
// 2. gets the HistoricalRates
// 3. computes the recommendation
func (h *Handler) computeRecommendation(currency string) (calculator.Signal, error) {
	startDate, endDate := date.GenerateStartAndEnd(DaysForRates)
	historicalRates, err := h.fx.GetHistoricalRates(currency, startDate, endDate)
	if err != nil || historicalRates == nil {
		return "", err
	}
	return h.ce.Recommend(historicalRates.RatesList), nil
}

func extractTargetRate(l *model.LatestRate) (float64, error) {
	if l == nil {
		return 0, errors.New("can't extract currency")
	}

	rates := l.Rates
	if r, ok := rates[calculator.EUR]; ok {
		return r, nil
	}
	return 0, errors.New("can't extract currency")
}
