package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"

	"FichainCore/crypto"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"

	"FichainBridge/config"
	"FichainBridge/models"   // Your models package
	"FichainBridge/requests" // Your requests DTO package
	"FichainBridge/scanner"
)

// --- Hard-coded values for development ONLY ---
// !!! DO NOT USE IN PRODUCTION !!!
const devPassword = "a-very-strong-and-secret-password-for-dev"

// Salt should ideally be unique per user, but for a hard-coded password,
// a hard-coded salt is consistent with the dev-only pattern.
var devSalt = []byte{0xc8, 0x28, 0x4c, 0xfd, 0x2, 0x55, 0xaa, 0x7}

type DepositWalletController struct {
	DB      *gorm.DB
	Scanner *scanner.BlockScannerService
}

// NewDepositWalletController creates and returns a new TransactionController.
func NewDepositWalletController(
	db *gorm.DB,
	scanner *scanner.BlockScannerService,
) *DepositWalletController {
	return &DepositWalletController{
		DB:      db,
		Scanner: scanner,
	}
}

// GET /api/v1/deposit-wallet/:address/:tokenName
func (ctrl *DepositWalletController) GetDepositWallet(c *gin.Context) {
	// 1. Get and Validate Input from Path Parameters
	fichainAddressHex := c.Param("address")
	if fichainAddressHex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address path parameter is required"})
		return
	}
	tokenName := c.Param("tokenName")
	if tokenName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tokenName path parameter is required"})
		return
	}
	if _, ok := config.GetConfig().TokenMap[tokenName]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tokenName"})
		return
	}

	// Clean up input "0x" prefix and decode from hex to bytes
	fichainAddressBytes, err := hex.DecodeString(strings.TrimPrefix(fichainAddressHex, "0x"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address format"})
		return
	}

	// 2. Find Existing Wallet in Database using BOTH fields
	var wallet models.DepositWallet
	result := ctrl.DB.Where(&models.DepositWallet{
		FichainAddress: fichainAddressBytes,
		TokenName:      tokenName,
	}).First(&wallet)

	// Check the result of the query
	if result.Error == nil {
		// --- WALLET FOUND ---
		// We found an existing wallet. Convert it to the response DTO.
		response := requests.ToDepositWalletResponse(&wallet)
		c.JSON(http.StatusOK, response)
		return
	}

	// If the error is anything other than "record not found", it's an internal server error.
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "database error while searching for wallet"},
		)
		return
	}

	// --- WALLET NOT FOUND, PROCEED TO CREATE ---
	// 3. Generate New BSC Key Pair
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new private key"})
		return
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	bscAddress := crypto.PubkeyToAddress(*publicKey)

	// 4. Encrypt the Private Key
	encryptedPrivateKey, err := encryptData(privateKeyBytes, devPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt private key"})
		return
	}

	// 5. Create and Save the New Wallet Record
	newWallet := &models.DepositWallet{
		FichainAddress:      fichainAddressBytes,
		TokenName:           tokenName, // <-- ADD THE TOKEN NAME
		Address:             bscAddress.Bytes(),
		EncryptedPrivateKey: encryptedPrivateKey,
		Status:              models.StatusActive,
	}

	createResult := ctrl.DB.Create(newWallet)
	if createResult.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to save new wallet to database"},
		)
		return
	}

	ctrl.Scanner.AddWalletToCache(newWallet)
	// 6. Return the newly created wallet in the response
	response := requests.ToDepositWalletResponse(newWallet)
	c.JSON(http.StatusCreated, response)
}

// encryptData encrypts plaintext data using AES-256-GCM.
// The key is derived from the password using scrypt.
// The output format is [nonce][ciphertext][tag].
func encryptData(plaintext []byte, password string) ([]byte, error) {
	// 1. Derive a 32-byte key from the password and salt using scrypt
	key, err := scrypt.Key([]byte(password), devSalt, 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	// 2. Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 3. Create a new GCM block cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 4. Create a nonce. GCM's nonce size is standard.
	// The nonce must be unique for each encryption with the same key.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 5. Encrypt the data. The nonce is prepended to the ciphertext.
	// gcm.Seal authentication tag is appended to the end of the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decryptData decrypts data encrypted by encryptData.
func decryptData(ciphertext []byte, password string) ([]byte, error) {
	// 1. Derive the key from the password and salt, same as in encryption
	key, err := scrypt.Key([]byte(password), devSalt, 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	// 2. Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 3. Create a new GCM block cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 4. Extract the nonce from the ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 5. Decrypt the data. gcm.Open will also verify the authentication tag.
	// If the tag is invalid (data was tampered with), it will return an error.
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
