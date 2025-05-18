package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

var userPool = sync.Pool{
	New: func() interface{} {
		return new(User)
	},
}

func getUsers(r io.Reader) (users []User, err error) {
	users = make([]User, 0, 100_000)

	scanner := bufio.NewScanner(r)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	buffer := make([]byte, 0, 1024)
	scanner.Buffer(buffer, 1024*1024)

	for scanner.Scan() {
		user := userPool.Get().(*User)
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return
		}
		users = append(users, *user)
		userPool.Put(user)
	}

	err = scanner.Err()
	return
}

func countDomains(u []User, domain string) (DomainStat, error) {
	domainName := "." + domain
	result := make(DomainStat, 500)

	for _, user := range u {
		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, domainName) {
			_, fullDomainName, _ := strings.Cut(email, "@")
			result[fullDomainName]++
		}
	}

	return result, nil
}
