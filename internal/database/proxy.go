package database

import (
	"time"

	"github.com/highbuyer/proxypool/log"
	"github.com/highbuyer/proxypool/pkg/proxy"
	"github.com/highbuyer/proxypool/redis"
	"gorm.io/gorm"
)

// 设置数据库字段，表名为默认为type名的复数。相比于原作者，不使用软删除特性
type Proxy struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	proxy.Base
	Link       string
	Identifier string `gorm:"unique"`
}

func InitTables() {
	if DB == nil {
		err := connect()
		if err != nil {
			return
		}
	}

	// 连接 Redis 数据库并进行测试
	redisConn, err2 := redis.ConnectRedis("localhost:6379", "")
	if err2 != nil {
		log.Errorln("\n\t\t[db/proxy.go] failed to connect to Redis database")
		panic(err)
	}

	_, pingErr := redisConn.Ping().Result()
	if pingErr != nil {
		log.Errorln("\n\t\t[db/proxy.go] failed to ping the Redis database")
		panic(pingErr)
	}

	err3 := DB.AutoMigrate(&Proxy{})
	if err3 != nil {
		log.Errorln("\n\t\t[db/proxy.go] database migration failed")
		panic(err3)
	}
}
func SaveProxyList(pl proxy.ProxyList) {
	if DB == nil {
		return
	}

	DB.Transaction(func(tx *gorm.DB) error {
		// Set All Usable to false
		if err := DB.Model(&Proxy{}).Where("useable = ?", true).Update("useable", false).Error; err != nil {
			log.Warnln("database: Reset useable to false failed: %s", err.Error())
		}
		// Create or Update proxies
		for i := 0; i < pl.Len(); i++ {
			p := Proxy{
				Base:       *pl[i].BaseInfo(),
				Link:       pl[i].Link(),
				Identifier: pl[i].Identifier(),
			}
			p.Useable = true
			if err := DB.Create(&p).Error; err != nil {
				// Update with Identifier
				if uperr := DB.Model(&Proxy{}).Where("identifier = ?", p.Identifier).Updates(&Proxy{
					Base: proxy.Base{Useable: true, Name: p.Name},
				}).Error; uperr != nil {
					log.Warnln("\n\t\tdatabase: Update failed:"+
						"\n\t\tdatabase: When Created item: %s"+
						"\n\t\tdatabase: When Updated item: %s", err.Error(), uperr.Error())
				}
			}
		}
		log.Infoln("database:Succeeded in saving the list of proxies")
		return nil
	})
}

// Get a proxy list consists of all proxies in database
func GetAllProxies() *proxy.ProxyList {
	proxies = make(proxy.ProxyList, 0)
	if DB == nil {
		return &proxies
	}

	proxiesDB := make([]Proxy, 0)
	DB.Select("link").Find(&proxiesDB)

	for _, proxyDB := range proxiesDB {
		if proxiesDB != nil {
			p, err := proxy.ParseProxyFromLink(proxyDB.Link)
			if err == nil && p != nil {
				p.SetUseable(false)
				*proxies = append(*proxies, p)
			}
		}
	}
	return proxies
}

// Clear proxies unusable more than 1 week
func ClearOldItems() {
	if DB == nil {
		return
	}

	lastWeek := time.Now().Add(-time.Hour * 24 * 7)
	if err := DB.Where("updated_at < ? AND useable = ?", lastWeek, false).Delete(&Proxy{}); err != nil {
		var count int64
		DB.Model(&Proxy{}).Where("updated_at < ? AND useable = ?", lastWeek, false).Count(&count)
		if count == 0 {
			log.Infoln("database: Nothing old to sweep") // TODO always this line?
		} else {
			log.Warnln("database: Delete old item failed: %s", err.Error())
		}
		return
	}
	log.Infoln("database: Succeeded in clearing old and unusable items")

}
