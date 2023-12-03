package services

import (
	"errors"
	"github.com/google/uuid"
	"gopay/internal/exts/db"
	"gopay/internal/models"
	"gorm.io/gorm"
	"time"
)

func CreateProductItems(productItems []models.ProductItem) error {
	if err := db.DB.Create(&productItems).Error; err != nil {
		return err
	}

	//更新库存数据
	if len(productItems) > 0 {
		if err := UpdateProductInStockCount([]uuid.UUID{productItems[0].ProductID}); err != nil {
			return err
		}
	}
	return nil
}

func DeleteProductItems(productItemIDs []uuid.UUID) error {
	result := db.DB.Model(&models.ProductItem{}).Delete("id in ?", productItemIDs)
	if result.RowsAffected == 0 {
		return errors.New("not_found")
	}
	return nil
}
func ClearExpireProductItem() error {
	query := db.DB.Model(&models.ProductItem{}).Where("status = 0 and end_lock_time < ?", time.Now().Unix()).Session(&gorm.Session{})
	var toUpdateInStockCountProductIDs []uuid.UUID
	if result := query.Pluck("product_id", &toUpdateInStockCountProductIDs); result.Error != nil {
		return errors.New("获取待更新商品失败")
	}
	if len(toUpdateInStockCountProductIDs) == 0 {
		return nil
	}

	// 解锁商品项目,取消绑定订单
	if result := query.Updates(map[string]interface{}{
		"status":        1,
		"order_id":      gorm.Expr("NULL"),
		"end_lock_time": gorm.Expr("NULL"),
	}); result.Error != nil {
		return errors.New("解锁商品项目失败")
	}

	if err := UpdateProductInStockCount(toUpdateInStockCountProductIDs); err != nil {
		return errors.New("更新商品库存失败")
	}

	return nil
}
