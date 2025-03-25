package server

import (
	"TarantoolKV/generated"
	"TarantoolKV/internal/application/core/domain"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var _ generated.ServerInterface = (*httpServer)(nil)

type application interface {
	Create(ctx context.Context, entity domain.Entity) error
	Update(ctx context.Context, entity domain.Entity) error
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (domain.Entity, error)
}
type httpServer struct {
	app application
}

func SetupHTTPServer(app application) *http.Server {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\" \"body size\": %d\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
			param.Request.ContentLength,
		)
	}))
	server := httpServer{
		app: app,
	}
	generated.RegisterHandlers(r, server)
	s := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}
	return s
}

func (h httpServer) PostKv(c *gin.Context) {
	var (
		data  map[string]interface{}
		key   string
		value map[string]interface{}
		ok    bool
	)
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if key, ok = data["key"].(string); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}
	if value, ok = data["value"].(map[string]interface{}); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "value is required"})
		return
	}
	err = h.app.Create(context.Background(), domain.Entity{Key: key, Value: value})
	if errors.Is(domain.ErrKeyExists, err) {
		c.JSON(http.StatusConflict, gin.H{"error": "key already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, value)
}

func (h httpServer) DeleteKvId(c *gin.Context, id string) {
	err := h.app.Delete(context.Background(), id)
	if errors.Is(domain.ErrKeyNotFound, err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (h httpServer) GetKvId(c *gin.Context, id string) {
	entity, err := h.app.Get(context.Background(), id)
	if errors.Is(domain.ErrKeyNotFound, err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": entity.Value})
	return
}

func (h httpServer) PutKvId(c *gin.Context, id string) {
	var (
		data  map[string]interface{}
		value map[string]interface{}
		ok    bool
	)
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if value, ok = data["value"].(map[string]interface{}); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "value is required"})
		return
	}
	err = h.app.Update(context.Background(), domain.Entity{Key: id, Value: value})
	if errors.Is(domain.ErrKeyNotFound, err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
