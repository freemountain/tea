// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
)

// AddLogin add login to config ( global var & file)
func AddLogin(name, token, user, passwd, sshKey, giteaURL string, insecure bool) error {
	// checks ...
	// ... if we have a url
	if len(giteaURL) == 0 {
		log.Fatal("You have to input Gitea server URL")
	}

	err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range Config.Logins {
		// ... if there already exist a login with same name
		if strings.ToLower(l.Name) == strings.ToLower(name) {
			return fmt.Errorf("login name '%s' has already been used", l.Name)
		}
		// ... if we already use this token
		if l.Token == token {
			return fmt.Errorf("token already been used, delete login '%s' first", l.Name)
		}
	}

	// .. if we have enough information to authenticate
	if len(token) == 0 && (len(user)+len(passwd)) == 0 {
		log.Fatal("No token set")
	} else if len(user) != 0 && len(passwd) == 0 {
		log.Fatal("No password set")
	} else if len(user) == 0 && len(passwd) != 0 {
		log.Fatal("No user set")
	}

	// Normalize URL
	serverURL, err := utils.NormalizeURL(giteaURL)
	if err != nil {
		log.Fatal("Unable to parse URL", err)
	}

	login := Login{
		Name:     name,
		URL:      serverURL.String(),
		Token:    token,
		Insecure: insecure,
		SSHKey:   sshKey,
		Created:  time.Now().Unix(),
	}

	if len(token) == 0 {
		login.Token, err = GenerateToken(login.Client(), user, passwd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Verify if authentication works and get user info
	u, _, err := login.Client().GetMyUserInfo()
	if err != nil {
		log.Fatal(err)
	}
	login.User = u.UserName

	if len(login.Name) == 0 {
		login.Name, err = GenerateLoginName(giteaURL, login.User)
		if err != nil {
			log.Fatal(err)
		}
	}

	// we do not have a method to get SSH config from api,
	// so we just use the hostname
	login.SSHHost = serverURL.Hostname()

	if len(sshKey) == 0 {
		login.SSHKey, err = login.FindSSHKey()
		if err != nil {
			fmt.Printf("Warning: problem while finding a SSH key: %s\n", err)
		}
	}

	// save login to global var
	Config.Logins = append(Config.Logins, login)

	// save login to config file
	err = SaveConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Login as %s on %s successful. Added this login as %s\n", login.User, login.URL, login.Name)

	return nil
}

// GenerateToken creates a new token when given BasicAuth credentials
func GenerateToken(client *gitea.Client, user, pass string) (string, error) {
	gitea.SetBasicAuth(user, pass)(client)

	host, _ := os.Hostname()
	tl, _, err := client.ListAccessTokens(gitea.ListAccessTokensOptions{})
	if err != nil {
		return "", err
	}
	tokenName := host + "-tea"

	for i := range tl {
		if tl[i].Name == tokenName {
			tokenName += time.Now().Format("2006-01-02_15-04-05")
			break
		}
	}

	t, _, err := client.CreateAccessToken(gitea.CreateAccessTokenOption{Name: tokenName})
	return t.Token, err
}
