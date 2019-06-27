package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
)

// FFXIVServers contains server info
type FFXIVServers struct {
	ServerNames        []string
	ServerCategory     []string
	ServerAvailability []bool
}

func main() {
	n := make(chan bool, 1)

	go runScrapper(n)

	d, _ := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	d.AddHandler(ready)

	d.Open()

	go func() {
		payload := <-n
		if payload {
			for i := 0; i < 5; i++ {
				d.ChannelMessageSend(os.Getenv("CHANNEL_ID"), "FAERIE IS OPEN @everyone")
				time.Sleep(5 * time.Second)
			}
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	d.Close()

}

func skipTags(num int, t *html.Tokenizer) {
	for i := 0; i < num; i++ {
		t.Next()
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Monitoring Lodestone...")
}

func runScrapper(not chan bool) {
	for {
		resp, _ := http.Get("https://na.finalfantasyxiv.com/lodestone/worldstatus/")

		b := resp.Body
		z := html.NewTokenizer(b)

		servers := &FFXIVServers{}

		done := false
		for done == false {
			nextTag := z.Next()
			switch {
			case nextTag == html.ErrorToken:
				done = true
			case nextTag == html.StartTagToken:
				t := z.Token()
				for _, g := range t.Attr {
					if g.Key == "class" {
						switch g.Val {
						case "world-list__world_name":
							skipTags(3, z)
							servers.ServerNames = append(servers.ServerNames, z.Token().Data)
						case "world-list__world_category":
							skipTags(3, z)
							servers.ServerCategory = append(servers.ServerCategory, z.Token().Data)
						case "world-ic__unavailable js__tooltip":
							servers.ServerAvailability = append(servers.ServerAvailability, false)
						case "world-ic__available js__tooltip":
							servers.ServerAvailability = append(servers.ServerAvailability, true)
						}
					}
				}
			}
		}

		for k, v := range servers.ServerNames {
			if v != "Faerie" {
				continue
			}

			if servers.ServerCategory[k] != "Congested" || servers.ServerAvailability[k] == false {
				not <- true
			}
		}

		// A little bit of cleanup
		servers = &FFXIVServers{}
		b.Close()
		time.Sleep(time.Minute)
	}
}
