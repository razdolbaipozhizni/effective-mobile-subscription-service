package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/model"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(svc service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: svc}
}

// Create godoc
// @Summary Создать новую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sub.StartDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date is required"})
		return
	}

	if err := h.service.Create(&sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// GetAll godoc
// @Summary Получить список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c *gin.Context) {
	subs, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// GetByID godoc
// @Summary Получить подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} model.Subscription
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	sub, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Update godoc
// @Summary Обновить данные подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body model.Subscription true "Новые данные"
// @Success 200 {object} model.Subscription
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(uint(id), &sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete godoc
// @Summary Удалить подписку
// @Tags subscriptions
// @Param id path int true "ID подписки"
// @Success 204 "No Content"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetTotalCost godoc
// @Summary Подсчёт суммарной стоимости подписок за период
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Param start_period query string true "Начало периода (MM-YYYY)"
// @Param end_period query string true "Конец периода (MM-YYYY)"
// @Success 200 {object} object{total=int}
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	startStr := c.Query("start_period")
	endStr := c.Query("end_period")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_period and end_period are required"})
		return
	}

	startPeriod, err := parsePeriod(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_period format (expected MM-YYYY)"})
		return
	}

	endPeriod, err := parsePeriod(endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_period format (expected MM-YYYY)"})
		return
	}

	var uIDPtr *uuid.UUID
	if userIDStr != "" {
		parsedUID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format (UUID required)"})
			return
		}
		uIDPtr = &parsedUID
	}

	var svcNamePtr *string
	if serviceName != "" {
		svcNamePtr = &serviceName
	}

	total, err := h.service.GetTotalCost(uIDPtr, svcNamePtr, startPeriod, endPeriod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}

func parsePeriod(period string) (time.Time, error) {
	return time.Parse("01-2006", period)
}