package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type CreateURLRequest struct {
	URL string
}

// will load this from env later on
var base62Alphabet = "QWERTYPOIUHBqa7412589630zxswedcvfrtgbnhyujmklpoiGVFCDXSZAJNMKL"

func base62encode(value int) string {

	res := ""

	for value > 0 {
		rem := value % 62
		value /= 62

		res = string(base62Alphabet[rem]) + res
	}

	return res
}

func getTicket() string {
	//GET TICKET FROM TICKET SERVER
	ticket := 1000

	encoded := base62encode(ticket)
	return encoded
}

func GetURL(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		ctx := context.Background()

		longURL, err := rdb.Get(ctx, shortURL).Result()

		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "URL does not exist",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve long URL",
			})
			return
		}
		c.Redirect(http.StatusFound, longURL)
	}
}

func CreateURL(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateURLRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Request",
			})
			return
		}

		longURL := req.URL

		shortURL := getTicket()

		ctx := context.Background()
		if err := rdb.Set(ctx, shortURL, longURL, 30*24*time.Hour).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create short URL",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"shortURL": shortURL,
		})
	}
}

//1month

//30 * 24 * 60 * 60 = 2592000
