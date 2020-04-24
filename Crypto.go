package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func CheckLicences(s *discordgo.Session) {
	directory, _ := filepath.Abs("./licences")

	d, err := os.Open(directory)
	if err != nil {
		PrintRedln("You don't have a folder named 'licences', cannot check for authenticity.")
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		PrintRedln("You don't have any files in the folder 'licences', cannot check for authenticity.")
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".lic" {
				data, err := ioutil.ReadFile(directory + "/" + file.Name()) // For read access.
				if err != nil {
					PrintRedln("Couldn't open : '" + file.Name() + "'")
				}
				result := Decrypt([]byte(data), "854E4DCDDCBA9DDA0A32139B36A14953D7269EC0346235E0D6DBF4E916AFFE8A")
				if strings.Contains(string(result), s.State.User.ID) {
					PrintMagentaln("Your licence has been validated, have fun !")
					Config.IsAllowedToUse = true
				}
			}
		}
	}
	if !Config.IsAllowedToUse {
		PrintRedln("If you didn't buy the bot, consider buying it, if you did, send a DM to Yewolf to solve the issue, or put your licence file in the 'licences' folder.")
	}
}
