package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// 使用bufio.NewReader從標準輸入創建一個新的讀取器
	reader := bufio.NewReader(os.Stdin)
	var err error
	var key1 string
	var key2 []byte
	for {
		fmt.Print("請輸入key1(16碼數字): ")
		key1, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("讀取輸入時發生錯誤:", err)
			return
		}
		key1 = strings.Replace(key1, "\n", "", -1)
		if _, err = strconv.Atoi(key1); err != nil {
			fmt.Println("請輸入數字")
			continue
		}
		count := len(key1)
		if len(key1) == 16 {
			break
		} else if count > 16 {
			fmt.Println("字數太多請再次輸入")
		} else {
			fmt.Println("字數不夠請再次輸入")
		}

	}
	for {
		fmt.Print("請輸入key2(12個文字或數字): ")
		key2, err = reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("讀取輸入時發生錯誤:", err)
			return
		}

		if len(key2) >= 12 {
			key2 = key2[0:12]
			fmt.Println(key2)
			break
		}
		fmt.Println("字數不夠請再次輸入")
	}

	repeatChan := make(chan struct{})
	defer close(repeatChan)
	// 使用reader.ReadString讀取輸入的文字，直到遇到換行符為止
	// 返回讀取到的文字和可能出現的錯誤
	for {
		fmt.Print("請輸入欲加密文字: ")

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("讀取輸入時發生錯誤:", err)
			return
		}

		go func() {
			// 原始資料
			plaintext := []byte(text)

			// 創建一個新的AES加密區塊，使用指定的金鑰
			block, err := aes.NewCipher([]byte(key1))
			if err != nil {
				panic(err)
			}

			// 創建一個Galois Counter Mode (GCM)的AES加密區塊，使用上面創建的AES加密區塊
			// 在這個示例中，我們使用了一個固定的nonce（並不安全，請勿在實際應用中使用）
			// 在實際應用中，nonce應該是一個隨機生成的位元組序列，長度為12位元組
			nonce := key2
			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				panic(err)
			}

			// 加密
			ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

			fmt.Printf("加密後的密文: %x\n", ciphertext)

			// 解密
			decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
			if err != nil {
				panic(err)
			}

			fmt.Printf("解密後的明文: %s\n", decrypted)
			repeatChan <- struct{}{}
		}()
		<-repeatChan
	}

}
