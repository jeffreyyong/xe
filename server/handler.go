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
)

type Handler struct {
	fx client.Forex
	ce calculator.Engine
}

func NewHandler(forex client.Forex, calculator calculator.Engine) *Handler {
	return &Handler{
		fx: forex,
		ce: calculator,
	}
}

var noRouteFoundFunc = func(c *gin.Context) {
	c.JSON(http.StatusNotFound, &model.ConvertResp{Error: model.ErrRouteNotFound})
}

func SetupAPIHandler(h *Handler) *gin.Engine {
	r := gin.Default()
	r.NoRoute(noRouteFoundFunc)
	r.GET(model.ConvertEndpoint, h.Convert)
	return r
}

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

	// Get latest rate
	latestRate, err := h.fx.GetLatestRate(currency)
	if err != nil {
		return http.StatusInternalServerError, &model.ConvertResp{Error: model.ErrConvert}, err
	}

	targetRate, err := extractTargetRate(latestRate)
	if err != nil {
		return http.StatusInternalServerError, &model.ConvertResp{Error: model.ErrConvert}, err
	}

	startDate, endDate := date.GenerateStartAndEnd(7)
	historicalRates, err := h.fx.GetHistoricalRates(currency, startDate, endDate)
	if err != nil {
		return http.StatusInternalServerError, &model.ConvertResp{Error: model.ErrConvert}, err
	}
	recommendation := h.ce.Recommend(historicalRates.RatesList)

	convertResp := &model.ConvertResp{
		From:           currency,
		To:             calculator.EUR,
		Rate:           targetRate,
		Recommendation: string(recommendation),
	}
	return http.StatusOK, convertResp, nil
}

func extractTargetRate(l *client.LatestRate) (float64, error) {
	if l == nil {
		return 0, errors.New("can't extract currency")
	}

	rates := l.Rates
	if r, ok := rates[calculator.EUR]; ok {
		return r, nil
	}
	return 0, errors.New("can't extract currency")
}
