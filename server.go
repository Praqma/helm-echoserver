package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	banner = " _    _      _             ______     _           _____\n" +
		"| |  | |    | |           |  ____|   | |         / ____|\n" +
		"| |__| | ___| |_ __ ___   | |__   ___| |__   ___| (___   ___ _ ____   _____ _ __\n" +
		"|  __  |/ _ | | '_ ` _ \\  |  __| / __| '_ \\ / _ \\\\___ \\ / _ | '__\\ \\ / / _ | '__|\n" +
		"| |  | |  __| | | | | | | | |___| (__| | | | (_) ____) |  __| |   \\ V |  __| |\n" +
		"|_|  |_|\\___|_|_| |_| |_| |______\\___|_| |_|\\___|_____/ \\___|_|    \\_/ \\___|_|"
	contentFile = "/tmp/content"
)

var p string

func main() {
	flag.StringVar(&p, "p", "8080", "the port to expose the server on.")

	flag.Parse()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/env", env)
	e.GET("/content", content)
	e.POST("/content", writeContent)
	// Start server
	e.Logger.Fatal(e.Start(":" + p))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, getHello()+getSpecificEnv("KUBERNETES")+getSpecificEnv("HELM")+getSpecificEnv("POD")+getSpecificEnv("CUSTOM")+getHeaders(c)+getResponse(c)+getHelp())
}

// Handler
func env(c echo.Context) error {
	return c.String(http.StatusOK, getEnv())
}

// Handler
func content(c echo.Context) error {
	return c.String(http.StatusOK, getContent())
}

// Handler
func writeContent(c echo.Context) error {
	m := c.FormValue("message")
	return c.String(http.StatusOK, appendContent(m))
}

func getHello() string {
	s := banner + "\nHello, I am Helm Echoserver, I echo some cool stuff about myself!\n\n"
	return s
}

func getSpecificEnv(q string) string {
	s := "\n\n -----> " + strings.ToUpper(q) + " \n\n"
	l := extractEnv(strings.ToUpper(q))
	if len(l) > 0 {
		for _, e := range l {
			pair := strings.Split(e, "=")
			s = s + pair[0] + "=" + pair[1] + "\n"
		}
	} else {
		s = s + "-no " + strings.ToLower(q) + " info- "
	}
	return s

}

func getHeaders(c echo.Context) string {
	s := "\n\n -----> HEADERS \n\n"
	s = s + " -> Request : " + c.Request().Proto + " " + c.Request().Method + " " + c.Request().URL.RequestURI() + "\n"
	s = s + " -> Host: " + c.Request().Host + "\n"
	s = s + " -> Remote IP: " + c.RealIP() + "\n"
	s = s + " -> Uri: " + c.Request().RequestURI + "\n"
	p := c.Request().URL.Path
	if p == "" {
		p = "/"
	}
	s = s + " -> Path: " + p + "\n"
	s = s + " -> Referer: " + c.Request().Referer() + "\n"
	s = s + " -> User Agent: " + c.Request().UserAgent() + "\n"
	l := c.Request().Header.Get(echo.HeaderContentLength)
	if l == "" {
		l = "0"
	}
	s = s + " -> Content Length: " + l + "\n"
	//s = s + " -> Content Length: " + c.Request().Header.Get(tag[7:]) + "\n"
	s = s + " -> Remote Address: " + c.Request().RemoteAddr + "\n"
	return s
}

func getResponse(c echo.Context) string {
	s := "\n\n -----> RESPONSE \n\n"
	if c.Response() != nil {
		s = s + " -> Status: " + strconv.Itoa(c.Response().Status) + "\n"
		s = s + " -> Content Length: " + strconv.FormatInt(c.Response().Size, 10) + "\n"
	} else {
		s = s + "-no response info-"
	}

	return s
}

func getHelp() string {
	s := "\n\n -----> Explore my other endpoints \n\n"
	s = s + " -> GET /env : prints all env variables where I am running.\n"
	s = s + " -> GET /content : reads the content of /tmp/content file where I am running.\n"
	s = s + " -> POST /content : writes (appends) to the content of /tmp/content file where I am running.\n"
	s = s + "e.g. curl -F \"message= a newely added content.\" http://localhost:" + p + "/content"
	return s
}

func getEnv() string {
	s := "\n -----> My environment |\n\n"
	envs := os.Environ()
	sort.Strings(envs)
	for _, e := range envs {
		pair := strings.Split(e, "=")
		s = s + pair[0] + "=" + pair[1] + "\n"
	}
	return s
}

func extractEnv(s string) []string {
	var result []string
	envs := os.Environ()
	sort.Strings(envs)
	for _, e := range envs {
		if strings.HasPrefix(e, s) {
			result = append(result, e)
		}
	}
	return result
}

func getContent() string {
	s, err := ioutil.ReadFile(contentFile)
	if err != nil {
		return "File " + contentFile + " does not exist."
	}
	return string(s)

}

func appendContent(m string) string {
	f, err := os.OpenFile(contentFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err.Error()
	}

	defer f.Close()

	if _, err = f.WriteString(m); err != nil {
		return err.Error()
	}
	return "Appended your message: [ " + m + " ] to " + contentFile
}
