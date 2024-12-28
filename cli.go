package gocli

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	goargs "github.com/budimanlai/go-args"
	goconfig "github.com/budimanlai/go-config"
)

type CliHandler = func(c *Cli)

type Cli struct {
	AppName    string
	Version    string
	IsShutdown bool

	Args        *goargs.Args
	Config      *goconfig.Config
	handler     map[string]CliHandler
	runtimePath string
	configFile  []string
}

type CliListen struct {
	OnLoop     func()
	OnShutdown func()
	TimeLoop   time.Duration
}

const (
	YYYYMMDDHHMMSS string = "2006-01-02 15:04:05"
)

func NewCli() *Cli {
	c := NewCliWithConfig(CliOptions{})

	return c
}

func NewCliWithConfig(config CliOptions) *Cli {
	c := &Cli{}

	c.Config = &goconfig.Config{}
	c.configFile = config.ConfigFile

	if config.AutoLoadConfig {
		if len(c.configFile) != 0 {
			e := c.Config.Open(config.ConfigFile...)

			if e != nil {
				panic(e)
			}
		}
	}

	c.Args = &goargs.Args{}
	c.Args.Parse()

	if config.RuntimePath == "" {
		config.RuntimePath = "runtime/"
	}
	c.runtimePath = config.RuntimePath

	c.handler = map[string]CliHandler{}
	c.addDefaultCommand()

	return c
}

func (c *Cli) LoadConfig() {
	if len(c.configFile) != 0 {
		e := c.Config.Open(c.configFile...)

		if e != nil {
			panic(e)
		}
	}
}

func (c *Cli) RuntimePath() string { return c.runtimePath }

func (c *Cli) addDefaultCommand() {
	c.AddCommandAndAlias("version", "v", func(c *Cli) {
		fmt.Println(c.AppName, "\nVersi", c.Version)
	})
}

func (c *Cli) AddCommandAndAlias(command string, alias string, handler CliHandler) {
	c.handler[command] = handler
	c.handler[alias] = handler
}

func (c *Cli) AddCommand(command string, handler CliHandler) {
	c.handler[command] = handler
}

func (c *Cli) Log(a ...interface{}) {
	now := time.Now()
	date := now.Format(YYYYMMDDHHMMSS)
	fmt.Print("[" + date + "] ")
	fmt.Println(a...)
}

func (c *Cli) Listen(handler CliListen) {
	// Channel untuk menangkap sinyal sistem
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP) // Tambahkan SIGHUP untuk reload

	var wg sync.WaitGroup // WaitGroup untuk sinkronisasi tugas
	var mutex sync.Mutex  // Mutex untuk status shutdown

	if handler.TimeLoop < 2 {
		handler.TimeLoop = 2 // Default minimum 2 detik
	}
	ticker := time.NewTicker(handler.TimeLoop * time.Second)
	defer ticker.Stop() // Pastikan ticker berhenti saat keluar

	// Goroutine untuk menangani sinyal
	go func() {
		for sig := range sigs {
			switch sig {
			case syscall.SIGHUP: // Reload konfigurasi
				c.Log("Sinyal reload konfigurasi diterima.")
				go func() {
					c.Log("Memuat ulang konfigurasi...")
					tempConfig := &goconfig.Config{}
					err := tempConfig.Open(c.configFile...)
					if err != nil {
						c.Log("Gagal memuat konfigurasi:", err)
						return
					}

					newConfig := *tempConfig // Buat salinan di luar mutex
					mutex.Lock()
					c.Config = &newConfig
					mutex.Unlock()
					c.Log("Konfigurasi berhasil dimuat ulang.")
				}()

			case syscall.SIGINT, syscall.SIGTERM: // Shutdown
				// Tandai shutdown
				mutex.Lock()
				c.IsShutdown = true
				mutex.Unlock()

				// Menjalankan handler OnShutdown jika ada
				if handler.OnShutdown != nil {
					wg.Add(1)

					// Eksekusi handler dalam goroutine terpisah
					go func() {
						defer wg.Done()
						defer func() {
							if r := recover(); r != nil {
								c.Log("Panic terdeteksi di OnShutdown:", r)
							}
						}()
						c.Log("Start shutdown...")
						handler.OnShutdown()
						c.Log("Shutdown done...")
					}()
				}

				wg.Wait()
				os.Exit(0) // Keluar setelah semua tugas selesai
			}
		}
	}()

	// Loop utama untuk menjalankan tugas
	c.IsShutdown = false
	if handler.OnLoop != nil {
		for range ticker.C { // Menggunakan for range
			mutex.Lock()
			if c.IsShutdown {
				mutex.Unlock()
				c.Log("Waiting shutdown")
				wg.Wait()
				return
			}
			mutex.Unlock()

			wg.Add(1)        // Sinkronisasi tugas
			handler.OnLoop() // Proses tugas (blocking)
			wg.Done()
		}
	}
}

func (c *Cli) Run() error {
	if h, e := c.handler[c.Args.Command]; e {
		h(c)
	} else {
		fmt.Print(fmt.Sprintf("Command '%s' not found", c.Args.Command) + "\n")
	}

	return nil
}
