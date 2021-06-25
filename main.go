package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/png"
	"log"
	"os"
)

var users []User
var hashLength int = 40

type User struct {
	Email     string
	Hash      string
	AvatarUrl string
}

func switchIndexDirection() func() int {
	idx := -1
	currentSign := "+"
	positiveSign := "+"
	negativeSign := "-"

	return func() int {
		if idx == 0 || idx == -1 {
			currentSign = positiveSign
		} else if idx == hashLength-1 {
			currentSign = negativeSign
		}
		if currentSign == positiveSign {
			idx++
			return idx
		}
		idx--
		return idx
	}
}
func generateHash(s string) string {

	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}
func generateAvatar(s string) *image.RGBA {
	size := 64
	sideBlocks := 3
	scale := size / sideBlocks
	img := image.NewRGBA(image.Rect(0, 0, sideBlocks*scale, sideBlocks*scale))

	idx := switchIndexDirection()
	for x := 0; x < sideBlocks; x++ {
		for y := 0; y < sideBlocks; y++ {
			i := idx()
			col := palette.Plan9[s[i]]
			startPoint := image.Point{x * scale, y * scale}
			endPoint := image.Point{x*scale + scale, y*scale + scale}
			rectangle := image.Rectangle{startPoint, endPoint}
			draw.Draw(img, rectangle, &image.Uniform{col}, image.Point{}, draw.Src)
		}
	}

	return img
}
func saveImage(img *image.RGBA, filename string) {

	f, err := os.Create(fmt.Sprintf("%s.png", filename))
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
func main() {
	// get user info & image size
	userEmail := flag.String("email", "email@default.com", "user email")

	flag.Parse()
	// generate hash
	hash := generateHash(*userEmail)
	// save user
	users = append(users, User{Email: *userEmail, Hash: hash})
	fmt.Printf("%v\n", users)
	// generate avatar
	avatar := generateAvatar(hash)
	// save avatar to some cloud storage
	saveImage(avatar, hash)
}
