package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"
		"log"
		"sync"
		    "math/rand"

	. "gopkg.in/check.v1"

	// "../cmd/lb/bala"
)

const baseAddress = "http://localhost:8090"

type server struct {
	name string
	connCnt int
}

var (
	serversPool = []*server {
		&server {
			name: "server1:8080",
			connCnt: 0,
		},
		&server {
			name: "server2:8080",
			connCnt: 0,
		},
		&server {
			name: "server3:8080",
			connCnt: 0,
		},
	}
)

func findMin(serversPool []*server) string {
	serverName := ""
	connMin := int(^uint(0) >> 1) // max int
	for i := 0; i < 3; i++ {
		server := serversPool[i]

		if (*server).connCnt < connMin {
			serverName = (*server).name
			connMin = (*server).connCnt
		}
	}
	return serverName
}

var client = http.Client{
	Timeout: 3 * time.Second,
}

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBalancer(c *C) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup){

			// time.Tick(time.Duration(rand.Intn(5)) * time.Second)

			min := findMin(serversPool)

			resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
			if err != nil {
				c.Logf("Response error%s", err)
			}

			server := resp.Header.Get("lb-from")
			log.Printf("response from [%s]", server)

			go func(sname string) {
				for i := 0; i < 3; i++ {
					if (*serversPool[i]).name == sname {
						(*serversPool[i]).connCnt++
						time.Tick(2 * time.Second)
						(*serversPool[i]).connCnt--
					}
				}
			}(server)

			defer wg.Done()

			c.Assert(server, Equals, min)
		}(&wg)
	}
	wg.Wait()

	time.Tick(5 * time.Second)
}


func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка.

	for i := 0; i < b.N; i++ {
		b.Skip("...skipping BenchmarkBalancer test...")
		fmt.Sprintf("------> BenchmarkBalancer <--------")
  }

}
