package httpapi

import (
	"log"
	"pick-up/backend/internal/db"
	"pick-up/backend/internal/httpapi/resp"
	"pick-up/backend/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	OrderController struct{}
	OrderUri        struct {
		ID uint `uri:"orderId" binding:"required" valid:"int"`
	}
	CreateOrderInput struct {
		StoreId uint `form:"storeId" binding:"required" valid:"int"`
	}
	CreateOrderResponse struct {
		Order models.Order `json:"order"`
	}
	GetOrderStatusResponse struct {
		Status models.OrderStatus `json:"status"`
	}
	GetOrderResponse struct {
		Order models.Order `json:"order"`
	}
)

func (ctl OrderController) CreateOrder(c *gin.Context) {
	// input validation
	input := CreateOrderInput{}
	if err := c.Bind(&input); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, resp.BodyErrorResp)
		return
	}

	// get cart
	cart, err := models.GetCart(c.GetUint("uid"), input.StoreId)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(404, resp.DbReadErrorResp)
		return
	}

	itemsCount := db.Conn().Model(&cart).Association("Items").Count()
	if itemsCount == 0 {
		c.AbortWithStatusJSON(424, gin.H{
			"message": "購物車內必須有東西",
		})
		return
	}

	// create order
	order := models.Order{
		StoreId: cart.StoreId,
		UserId:  cart.UserId,
		CartId:  cart.ID,
		Price:   cart.CalculateTotalPrice(),
		Status:  models.OrderStatus(models.OrderStatusCreated),
	}
	err = db.Conn().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// This cart will be soft-deleted and archived for this order
		cart.DeletedAt = gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		}
		if err := tx.Save(cart).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, resp.DbWriteErrorResp)
		return
	}

	go simulateOrderProgress(&order)

	c.AbortWithStatusJSON(200, CreateOrderResponse{
		Order: order,
	})
}

func (ctl OrderController) GetOrderStatus(c *gin.Context) {
	// input validation
	uri := OrderUri{}
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, resp.UriErrorResp)
		return
	}

	// db operation
	order := models.Order{}
	if err := db.Conn().Find(&order, uri.ID).Error; err != nil {
		c.AbortWithStatusJSON(404, gin.H{
			"message": "指定店家不存在",
		})
		return
	}

	c.AbortWithStatusJSON(200, GetOrderStatusResponse{
		Status: order.Status,
	})
}

func (ctl OrderController) GetOrder(c *gin.Context) {
	// input validation
	uri := OrderUri{}
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, resp.UriErrorResp)
		return
	}

	// db operation
	order := models.Order{}
	if err := db.Conn().Preload("Store").Preload("Cart").First(&order, uri.ID).Error; err != nil {
		c.AbortWithStatusJSON(404, gin.H{
			"message": "指定訂單不存在",
		})
		return
	}

	c.AbortWithStatusJSON(200, GetOrderResponse{
		Order: order,
	})
}
