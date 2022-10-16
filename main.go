package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"z3ntl3/token-checker-layered/builder"

	request "github.com/parnurzeal/gorequest"
	"golang.org/x/sync/errgroup"
	"gopkg.in/ini.v1"
)

const (
	API string = "https://discord.com/api/v9/users/@me"
)

/*
* Programmed by Z3NTL3 - pix4.dev
 */
var (
	goods int = 0
)
var proxiee string

type blueprint_Client interface {
	startTokenCheck() error
}

type CLIENT struct {
	req   *request.SuperAgent
	token string
	file  *os.File
}

func (c CLIENT) startTokenCheck() error {
	token := c.token
	req := c.req

	resp, body, err := req.Proxy("http://"+proxiee).Get(API).Set("Authorization", token).Set("user-agent", "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36").Timeout(time.Second * 5).End()
	tunnel := strings.Split(proxiee, "@")[1]

	if err != nil {
		fmt.Println(err)
		return errors.New("\033[1m\033[0;97m[INFO] \033[31mCant Connect \033[1m\033[0;97m- Proxy Issue \033[0;97mFor info look above at the error log\033[0m")
	}
	defer resp.Body.Close()
	fmt.Println("\033[1m\033[0;97m[INFO] \033[0;95mTrying to tunnel on PROXY \033[1m\033[0;97m[ " + tunnel + " ]\033[0m")
	fmt.Println("\033[1m\033[0;97m[INFO] \033[0;32mTunnel success PROXIED! \033[1m\033[0;97m[ " + tunnel + " ]\033[0m\n")

	if resp.StatusCode == 200 {
		if strings.Contains(body, "username") && !strings.Contains(body, "error") {
			fmt.Println("\033[1m\033[0;97m[INFO] \033[32mGood Token: \033[1m\033[0;97m", token, "\033[0m")
			c.file.WriteString(token + "\n")
			goods += 1
		} else {
			fmt.Println("\033[1m\033[0;97m[INFO] \033[31mBad Token: \033[1m\033[0;97m", token, "\033[0m")
		}
	} else if !strings.Contains(body, "error") {
		fmt.Println("\033[1m\033[0;97m[INFO] \033[31mBad Token: \033[1m\033[0;97m", token, "\033[0m")
	} else {
		return errors.New("\033[1m\033[0;97m[INFO] \033[31mRate Limit \033[1m\033[0;97m- Change Proxy and then start the tool again!\033[0m")
	}
	return nil
}

func inputCheck(logo string, args []string) ([]string, error) {
	fmt.Printf(logo)
	if len(args) != 2 {
		fmt.Println(builder.Usage(), "\n")
		return nil, errors.New("Invalid usage")
	}
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Cannot find/or read the file!")
	}

	if len(string(content)) == 0 {
		fmt.Println("\n\033[1m\033[0;97m[INFO] \033[31mAdd tokens into your token file\033[0m")
		return nil, errors.New("Add tokens into your token file")
	}
	clear := strings.Trim(string(content), "\n")
	tokens := strings.Split(clear, "\n")

	fmt.Println()
	return tokens, nil
}

func main() {
	logo := builder.LogoBuild()
	args := os.Args

	tokens, err := inputCheck(logo, args)

	max_worker_count := runtime.NumCPU()
	free_cores := 3

	workers := new(errgroup.Group)
	workers.SetLimit(10000 * (max_worker_count - free_cores))

	f, err := os.OpenFile("goods.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	cfg, err := ini.Load("proxy.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	proxiee = cfg.Section("proxy").Key("http_proxy").String()

	ssl := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	for _, v := range tokens {
		req := request.New()
		req.Transport = ssl

		v = strings.TrimSpace(v)
		info := CLIENT{
			req:   req,
			token: v,
			file:  f,
		}
		workers.Go(func() error {
			return info.startTokenCheck()
		})
	}

	if err := workers.Wait(); err != nil {
		fmt.Println("\n", err)
	}
	fmt.Println("\n\033[1m\033[0;97m[INFO] \033[0;95mAmount Valid: \033[1m\033[0;97m" + strconv.Itoa(goods) + ", \033[38;5;96mSaved in file \033[0m[\033[38;5;97mgoods.txt\033[0m]\033[0m")
}
