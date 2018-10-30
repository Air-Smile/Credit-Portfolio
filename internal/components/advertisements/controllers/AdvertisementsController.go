package controllers

import (
	"github.com/a1ta1r/Credit-Portfolio/internal/codes"
	"github.com/a1ta1r/Credit-Portfolio/internal/components/advertisements/entities"
	"github.com/a1ta1r/Credit-Portfolio/internal/components/advertisements/storages"
	"github.com/a1ta1r/Credit-Portfolio/internal/components/roles"
	_ "github.com/a1ta1r/Credit-Portfolio/internal/specification/responses"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

type AdvertisementsController struct {
	advertiserStorage    storages.AdvertiserStorage
	advertisementStorage storages.AdvertisementStorage
	bannerStorage        storages.BannerStorage
	bannerPlaceStorage   storages.BannerPlaceStorage
}

func NewAdvertisementController(
	advertiserStorage storages.AdvertiserStorage,
	advertisementStorage storages.AdvertisementStorage,
	bannerStorage storages.BannerStorage,
	bannerPlaceStorage storages.BannerPlaceStorage,
) AdvertisementsController {
	return AdvertisementsController{
		advertiserStorage:    advertiserStorage,
		advertisementStorage: advertisementStorage,
		bannerStorage:        bannerStorage,
		bannerPlaceStorage:   bannerPlaceStorage,
	}
}

// @Summary Получить список всех рекламодателей
// @Description Метод возвращает список всех имеющихся в системе рекламодателей
// @Produce  json
// @Success 200 {object} responses.AllAdvertisers
// @Router /advertisers [get]
func (ac AdvertisementsController) GetAdvertisers(c *gin.Context) {
	var advertisers []entities.Advertiser
	advertisers, _ = ac.advertiserStorage.GetAdvertisers()
	for i := 0; i < len(advertisers); i++ {
		advertisers[i].Password = ""
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      "OK",
		"count":       len(advertisers),
		"advertisers": advertisers,
	})
}

// @Summary Получить рекламодателя по ID
// @Description Метод возвращает рекламодателя по его ID
// @Produce  json
// @Param id path int true "ID рекламодателя"
// @Success 200 {object} entities.Advertiser "{"advertiser": entities.Advertiser}"
// @Failure 404 "{"message": "resource not found"}"
// @Router /advertisers/{id} [get]
func (ac AdvertisementsController) GetAdvertiser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": codes.BadID})
		return
	}
	advertiser, err := ac.advertiserStorage.GetAdvertiser(uint(id))
	if advertiser.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": codes.ResNotFound})
		return
	}
	advertiser.Password = ""
	c.JSON(http.StatusOK, gin.H{"advertiser": advertiser})
}

// @Summary Добавить нового рекламодателя
// @Description Метод добавляет в систему нового рекламодателя с заданными параметрами
// @Accept json
// @Produce  json
// @Param advertiser body entities.Advertiser true "Данные о рекламодателе"
// @Success 200 {object} entities.Advertiser "{"advertiser": entities.Advertiser}"
// @Router /advertisers [post]
func (ac AdvertisementsController) AddAdvertiser(c *gin.Context) {
	var advertiser entities.Advertiser
	c.BindJSON(&advertiser)
	advertiser.Role = roles.Ads
	advertiser.Password = advertiser.GetHashedPassword()
	err := ac.advertiserStorage.CreateAdvertiser(&advertiser)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": codes.ResourceExists})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"advertiser": advertiser})
}

func (ac AdvertisementsController) DeleteAdvertiser(c *gin.Context) {
	var advertiser entities.Advertiser
	c.BindJSON(&advertiser)
	ac.advertiserStorage.DeleteAdvertiser(advertiser)
	c.JSON(http.StatusOK, gin.H{"message": codes.ResDeleted})
}

func (ac AdvertisementsController) UpdateAdvertiser(c *gin.Context) {
	var advertiser entities.Advertiser
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": codes.BadID})
		return
	}
	advertiser, _ = ac.advertiserStorage.GetAdvertiser(uint(id))
	if advertiser.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": codes.ResNotFound})
		return
	}

	c.ShouldBindWith(&advertiser, binding.JSON)
	_ = ac.advertiserStorage.UpdateAdvertiser(advertiser)
	c.JSON(http.StatusOK, gin.H{
		"status":     "OK",
		"advertiser": advertiser,
	})
}

func (ac AdvertisementsController) GetAdvertisementsByAdvertiser(c *gin.Context) {
	var advertisements []entities.Advertisement
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": codes.BadID})
		return
	}
	advertisements, _ = ac.advertisementStorage.GetAdvertisementsByAdvertiser(uint(id))
	c.JSON(http.StatusOK, gin.H{
		"status":         "OK",
		"count":          len(advertisements),
		"advertisements": advertisements,
	})
}

func (ac AdvertisementsController) GetAdvertisements(c *gin.Context) {
	var advertisement []entities.Advertisement
	advertisement, _ = ac.advertisementStorage.GetAdvertisements()
	c.JSON(http.StatusOK, gin.H{
		"status":      "OK",
		"count":       len(advertisement),
		"advertisers": advertisement,
	})
}

func (ac AdvertisementsController) GetAdvertisement(c *gin.Context) {
	var advertisement entities.Advertiser
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": codes.BadID})
		return
	}
	advertisement, _ = ac.advertiserStorage.GetAdvertiser(uint(id))
	if advertisement.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": codes.ResNotFound})
		return
	}
	c.JSON(http.StatusOK, gin.H{"advertiser": advertisement})
}

func (ac AdvertisementsController) AddAdvertisement(c *gin.Context) {
	var advertiser entities.Advertiser
	c.BindJSON(&advertiser)
	err := ac.advertiserStorage.CreateAdvertiser(&advertiser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": codes.InternalError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"advertiser": advertiser})
}

func (ac AdvertisementsController) DeleteAdvertisement(c *gin.Context) {
	var advertiser entities.Advertiser
	c.BindJSON(&advertiser)
	ac.advertiserStorage.DeleteAdvertiser(advertiser)
	c.JSON(http.StatusOK, gin.H{"message": codes.ResDeleted})
}

func (ac AdvertisementsController) UpdateAdvertisement(c *gin.Context) {
	var advertiser entities.Advertiser
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": codes.BadID})
		return
	}
	advertiser, _ = ac.advertiserStorage.GetAdvertiser(uint(id))
	if advertiser.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": codes.ResNotFound})
		return
	}

	c.ShouldBindWith(&advertiser, binding.JSON)
	_ = ac.advertiserStorage.UpdateAdvertiser(advertiser)
	c.JSON(http.StatusOK, gin.H{
		"status":     "OK",
		"advertiser": advertiser,
	})
}

func (ac AdvertisementsController) GetBannersByAdvertisement(c *gin.Context) {
	var banners []entities.Banner
	id, err := strconv.ParseUint(c.Param("adsid"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": codes.BadID})
		return
	}
	banners, _ = ac.bannerStorage.GetBannersByAdvertisement(uint(id))
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"count":   len(banners),
		"banners": banners,
	})
}