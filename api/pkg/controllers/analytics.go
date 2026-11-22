package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/plutov/formulosity/api/pkg/http/response"
	surveyspkg "github.com/plutov/formulosity/api/pkg/surveys"
	"github.com/plutov/formulosity/api/pkg/types"
)

func (h *Handler) getSurveyAnalytics(c echo.Context) error {
	survey := c.Get("survey").(types.Survey)

	analytics, err := surveyspkg.GetSurveyAnalytics(h.Services, survey)
	if err != nil {
		h.Logger.Error("unable to get survey analytics", "err", err)
		return response.InternalErrorDefaultMsg(c)
	}

	return response.Ok(c, analytics)
}
