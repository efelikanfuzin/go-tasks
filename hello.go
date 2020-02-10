package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

type User struct {
	name  string
	email string
}

func (user *User) GenerateAvatar() string {
	return "avatr_link.png"
}

func (user *User) Hash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(user.name+user.email)))
}

func (user *User) generateAvatar() string {
	var color_multiplaer uint64 = 1
	hash := user.Hash()
	width := 500
	height := 500

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			slice_start := (x + y) % len(hash)

			if slice_start+4 > len(hash) {
				slice_start = 0
			}

			rgb_item, _ := strconv.ParseUint(hash[slice_start:slice_start+4], 16, 64)

			pixel_color := color.RGBA{
				uint8(rgb_item * (uint64(slice_start) * color_multiplaer)),
				uint8(rgb_item - (uint64(slice_start) * color_multiplaer)),
				uint8(rgb_item + (uint64(slice_start) * color_multiplaer)),
				0xff,
			}
			img.Set(x, y, pixel_color)
			color_multiplaer++
		}
	}
	f, _ := os.Create("image.png")
	png.Encode(f, img)

	return "image.png"
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')
	fmt.Printf("Enter your email")
	email, _ := reader.ReadString('\n')
	user := User{name, email}

	user.generateAvatar()
}
