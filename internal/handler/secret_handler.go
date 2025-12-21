package handler

import (
	"burn-text/crypto"
	"burn-text/internal"
	"burn-text/internal/config"
	"burn-text/storage"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateSecret(c *gin.Context) {
	var req internal.CreateSecretReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	//处理密码
	var pwdHash string
	if req.Password != "" {
		hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
			return
		}
		pwdHash = string(hashBytes)
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

	//构造存储对象
	dataObj := internal.SecretData{
		CipherText:   cipherText,
		PasswordHash: pwdHash,
	}
	dataJson, _ := json.Marshal(dataObj)

	//生成ID
	idBytes, _ := crypto.GenerateKey()
	id := hex.EncodeToString(idBytes[:8])

	//存入Redis
	ttl := time.Duration(config.GlobalConfig.Redis.TTLMinutes) * time.Minute
	err = storage.Save(id, string(dataJson), ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储失败"})
		return
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

	visitPassword := c.GetHeader("X-Access-Password")

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

	var dataObj internal.SecretData
	if err := json.Unmarshal([]byte(cipherText), &dataObj); err != nil {
		//兼容旧数据，如果解析失败直接使用cipherText
		dataObj.CipherText = cipherText
	}

	//验证密码
	if dataObj.PasswordHash != "" {
		if visitPassword == "" {
			//密码不对，数据塞回，但存在并发风险
			_ = storage.Save(id, cipherText, 10*time.Minute)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要密码访问"})
			return
		}

		//验证哈希
		err := bcrypt.CompareHashAndPassword([]byte(dataObj.PasswordHash), []byte(visitPassword))
		if err != nil {
			//密码不对，数据塞回，但存在并发风险
			_ = storage.Save(id, cipherText, 10*time.Minute)
			c.JSON(http.StatusForbidden, gin.H{"error": "访问密码错误"})
			return
		}
	}

	//解码Key
	keyBytes, _ := hex.DecodeString(keyString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密钥格式错误"})
		return
	}

	//解密文本
	keyBytes, _ = hex.DecodeString(keyString)
	plainText, err := crypto.Decrypt(dataObj.CipherText, keyBytes)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "密钥错误,解密失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "阅读成功",
		"data": plainText,
	})
}
