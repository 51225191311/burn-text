package handler

import (
	"burn-text/crypto"
	"burn-text/internal"
	"burn-text/internal/config"
	"burn-text/storage"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"

	"net/http"
	"time"
)

func CreateSecret(c *gin.Context) {
	var req internal.CreateSecretReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	//生成Key
	keyBytes, _ := crypto.GenerateKey()
	keyString := hex.EncodeToString(keyBytes)

	//加密文本
	cipherText, err := crypto.Encrypt(req.Text, keyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加密系统异常"})
		return
	}

	//生成ID
	idBytes, _ := crypto.GenerateKey()
	id := hex.EncodeToString(idBytes[:8])

	//存入Redis
	ttl := time.Duration(config.GlobalConfig.Redis.TTLMinutes) * time.Minute
	err = storage.Save(id, cipherText, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储失败"})
	}

	//生成URL
	port := config.GlobalConfig.Server.Port
	fullUrl := fmt.Sprintf("http://localhost:%s/view/%s?key=%s", port, id, keyString)

	c.JSON(http.StatusOK, internal.CreateSecretResp{
		ID:  id,
		Url: fullUrl,
	})
}

func GetSecret(c *gin.Context) {
	id := c.Param("id")
	keyString := c.Query("key")

	if keyString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少密钥"})
		return
	}

	//从Redis取并删
	cipherText, err := storage.GetAndDelete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "信息不存在或已阅后即焚"})
		return
	}

	//解码Key
	keyBytes, _ := hex.DecodeString(keyString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密钥格式错误"})
		return
	}

	//解密文本
	plainText, err := crypto.Decrypt(cipherText, keyBytes)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "密钥错误,解密失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "阅读成功",
		"data": plainText,
	})
}
