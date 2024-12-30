package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"slices"

	"errors"

	goSystemd "github.com/alirezasn3/go-systemd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	ID          string `json:"id" redis:"id"`
	FirstName   string `json:"firstName" redis:"firstName"`
	LastName    string `json:"lastName" redis:"lastName"`
	PhoneNumber string `json:"phoneNumber" redis:"phoneNumber"`
	Role        string `json:"-" redis:"role"`
	Password    string `json:"password" redis:"password"`
}

type Config struct {
	RedisURL       string   `json:"redisURL"`
	JWTSecret      string   `json:"jwtSecret"`
	Domain         string   `json:"domain"`
	AllowedOrigins []string `json:"allowedOrigins"`
}

var execDir string
var config Config
var redisClient *RedisWrapper

func runCommand(args []string, envVars []string, prefix string, shouldLog bool) error {
	cmd := exec.Command(args[0], args[1:]...)

	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	go func() {
		for outScanner.Scan() {
			text := outScanner.Text()
			if len(text) == 0 {
				continue
			}
			if shouldLog {
				log.Println(prefix + " " + text)
			}
		}
	}()

	errText := ""
	go func() {
		for errScanner.Scan() {
			text := errScanner.Text()
			if len(text) == 0 {
				continue
			}
			if len(errText) > 0 {
				errText += "\n"
			}
			errText += text
			if shouldLog {
				log.Println(prefix + " " + text)
			}
		}
	}()

	cmd.Wait()

	if err = outScanner.Err(); err != nil {
		return err
	}

	if len(errText) > 0 {
		return errors.New(errText)
	}

	if err = errScanner.Err(); err != nil {
		return err
	}

	return nil
}

func getExecDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	execDir := filepath.Dir(execPath)
	return execDir, nil
}

func init() {
	// get the executable directory
	var err error
	execDir, err = getExecDir()
	if err != nil {
		log.Println("failed to get executable directory")
		os.Exit(1)
	}

	// check for install and uninstall commands on linux
	if runtime.GOOS == "linux" {
		if slices.Contains(os.Args, "--install") {
			execPath, err := os.Executable()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = goSystemd.CreateService(&goSystemd.Service{Name: "echo-redis-svelte-template", ExecStart: execPath, Restart: "on-failure", RestartSec: "3s"})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				fmt.Println("echo-redis-svelte-template service created")
				os.Exit(0)
			}
		} else if slices.Contains(os.Args, "--uninstall") {
			err := goSystemd.DeleteService("echo-redis-svelte-template")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				fmt.Println("echo-redis-svelte-template service deleted")
				os.Exit(0)
			}
		}
	}

	bytes, err := os.ReadFile(filepath.Join(execDir, "config.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}
	log.Println("Loaded config from " + filepath.Join(execDir, "config.json"))
}

func main() {
	// connect to redis
	var err error
	redisClient, err = CreateRedisClient(config.RedisURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	} else {
		log.Println("Connected to redis")
	}

	fmt.Println(redisClient.KeyExists("test"))

	// check for dev mode and run frontend node server
	if slices.Contains(os.Args, "--dev") {
		// run the server in dev mode
		go func() {
			err := runCommand([]string{"npm", "run", "dev", "--prefix", path.Join(execDir, "frontend"), "--", "--port=1234"}, []string{}, "[frontend build]", true)
			if err != nil {
				log.Println(err)
				log.Println("failed to run frontend server")
			}
			log.Println("frontend server failed, exiting")
			os.Exit(1)
		}()
	} else {
		// build the frontend
		err := runCommand([]string{"npm", "run", "build", "--prefix", path.Join(execDir, "frontend")}, []string{}, "[frontend build]", true)
		if err != nil {
			log.Println("failed to build frontend")
		}
		log.Println("frontend build successful")

		// run the server
		go func() {
			err = runCommand([]string{"node", path.Join(execDir, "frontend", "build", "index.js")}, []string{"PORT=1234"}, "[frontend server]", true)
			if err != nil {
				log.Println(err)
				log.Println("failed to run frontend server")
			}
			log.Println("frontend server failed, exiting")
			os.Exit(1)
		}()
	}

	// create echo instance
	e := echo.New()

	// create routers
	apiRouter := e.Group("/api")
	frontendRouter := e.Group("/*")

	// parse frontend upstream
	frontendUpstream, err := url.Parse("http://localhost:1234")
	if err != nil {
		log.Println("failed to parse frontend upstream")
		os.Exit(1)
	}

	// proxy frontend requests to frontend server
	frontendRouter.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: frontendUpstream}})))

	// add cors middleware to api router
	apiRouter.Use(CORSMiddleware())

	// register api routes
	apiRouter.POST("/signup", Signup)
	apiRouter.POST("/signin", Signin)
	apiRouter.GET("/signout", Signout)
	apiRouter.GET("/account", Account, AuthenticationMiddleware)

	// register api not found handler
	apiRouter.RouteNotFound("/*", func(c echo.Context) error {
		return c.NoContent(404)
	})

	// start echo server
	e.Logger.Fatal(e.Start(":8080"))
}
