package gocli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	goargs "github.com/budimanlai/go-args"
	goconfig "github.com/budimanlai/go-config"
)

type CliHandler = func(c *Cli)

type Cli struct {
	IsShutdown bool

	Args        *goargs.Args
	Config      *goconfig.Config
	handler     map[string]CliHandler
	runtimePath string
	configFile  []string
	appName     string
	version     string
	wg          sync.WaitGroup // WaitGroup untuk sinkronisasi tugas
	mutex       sync.Mutex     // Mutex untuk status shutdown
	pidFile     string         // File untuk menyimpan PID
}

type CliRunLoop struct {
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

	c.appName = config.AppName
	c.version = config.Version

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

	c.pidFile = config.RuntimePath + c.appName + ".pid"

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
		fmt.Println(c.appName, "\nVersi", c.version)
	})
	c.AddCommand("start", func(c *Cli) {
		c.startDaemon()
	})
	c.AddCommand("stop", func(c *Cli) {
		c.stopDaemon()
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

func (c *Cli) listenSignal(handler CliRunLoop) {
	// Channel untuk menangkap sinyal sistem
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP) // Tambahkan SIGHUP untuk reload

	// Goroutine untuk menangani sinyal
	go func() {
		defer close(sigs) // Ensure the channel is closed
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
					c.mutex.Lock()
					c.Config = &newConfig
					c.mutex.Unlock()
					c.Log("Konfigurasi berhasil dimuat ulang.")
				}()

			case syscall.SIGINT, syscall.SIGTERM: // Shutdown
				// Tandai shutdown
				c.mutex.Lock()
				c.IsShutdown = true
				c.mutex.Unlock()

				// Menjalankan handler OnShutdown jika ada
				if handler.OnShutdown != nil {
					c.wg.Add(1)

					// Eksekusi handler dalam goroutine terpisah
					go func() {
						defer c.wg.Done()
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

				c.wg.Wait()
				signal.Stop(sigs) // Stop receiving signals
				os.Exit(0)        // Keluar setelah semua tugas selesai
			}
		}
	}()
}

func (c *Cli) RunLoop(handler CliRunLoop) {
	c.listenSignal(handler)

	if handler.TimeLoop < 2 {
		handler.TimeLoop = 2 * time.Second // Default minimum 2 detik
	}
	ticker := time.NewTicker(handler.TimeLoop)
	defer ticker.Stop() // Pastikan ticker berhenti saat keluar

	// Loop utama untuk menjalankan tugas
	c.IsShutdown = false
	if handler.OnLoop != nil {
		for range ticker.C { // Menggunakan for range
			c.mutex.Lock()
			if c.IsShutdown {
				c.mutex.Unlock()
				c.Log("Waiting shutdown")
				c.wg.Wait()
				return
			}
			c.mutex.Unlock()

			c.wg.Add(1) // Sinkronisasi tugas
			go func() {
				defer c.wg.Done() // Ensure Done is called even if OnLoop panics
				handler.OnLoop()  // Proses tugas (blocking)
			}()
		}
	}
}

func (c *Cli) Run() error {
	if h, exists := c.handler[c.Args.Command]; exists {
		h(c)
	} else {
		fmt.Printf("Command '%s' not found\n", c.Args.Command)
	}

	return nil
}

func (c *Cli) startDaemon() {
	logFile, err := os.OpenFile(c.runtimePath+"daemon.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer logFile.Close()

	cmd := exec.Command(os.Args[0], "run")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // Detach the process
	err = cmd.Start()
	if err != nil {
		fmt.Println("Failed to start daemon:", err)
		return
	}
	err = os.WriteFile(c.pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	if err != nil {
		fmt.Println("Failed to write PID file:", err)
		return
	}
	os.Exit(0)
}

func (c *Cli) stopDaemon() {
	pidData, err := os.ReadFile(c.pidFile)
	if err != nil {
		fmt.Println("Failed to read PID file:", err)
		return
	}
	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		fmt.Println("Invalid PID:", err)
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Failed to find process:", err)
		return
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Println("Failed to stop process:", err)
		return
	}
	os.Remove(c.pidFile)
}
