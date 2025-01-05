# Go-CLI

Go-CLI adalah library untuk membuat aplikasi command-line interface (CLI) dengan mudah menggunakan bahasa pemrograman Go. Library ini menyediakan berbagai fitur untuk mempermudah pembuatan CLI, seperti parsing argumen, sub-komando, dan bantuan otomatis.

## Kegunaan
Library ini berguna untuk:
- Membuat aplikasi CLI dengan cepat dan mudah.
- Mengelola argumen dan opsi dari command line.
- Menyediakan bantuan otomatis untuk pengguna CLI.

## Cara Menggunakan
1. Instalasi:
   ```sh
   go get github.com/username/go-cli
   ```

2. Contoh Penggunaan:
   Lihat: `example/main.go`

3. Membuat Aplikasi CLI:
   ```go
   package main

   import (
       "fmt"
       "github.com/username/go-cli"
   )

   func main() {
       app := gocli.NewApp()
       app.Name = "MyApp"
       app.Usage = "Ini adalah aplikasi CLI saya"

       app.Action = func(c *gocli.Context) error {
           fmt.Println("Hello, CLI!")
           return nil
       }

       app.Run(os.Args)
   }
   ```

4. Menambahkan Sub-komando:
   ```go
   app.Commands = []gocli.Command{
       {
           Name:    "greet",
           Aliases: []string{"g"},
           Usage:   "Menampilkan salam",
           Action: func(c *gocli.Context) error {
               fmt.Println("Hello, World!")
               return nil
           },
       },
   }
   ```

Untuk informasi lebih lanjut, silakan merujuk ke dokumentasi lengkap di [link dokumentasi].