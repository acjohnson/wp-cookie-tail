package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hpcloud/tail"
)

var ctx = context.Background()

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	var cookie_log string
	cookie_prefix := "wordpress_logged_in_"

	wp_log := getEnv("WP_LOG", "/tmp/wp_cookie_log")
	redis_host := getEnv("REDIS_HOST", "localhost")
	redis_port := getEnv("REDIS_PORT", "6379")
	redis_ttl, err := strconv.Atoi(getEnv("REDIS_TTL", "1209600"))
	if err != nil {
		log.Fatal(err)
	}

	redis_ttl_time := time.Duration(redis_ttl) * time.Second

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host + ":" + redis_port,
		Password: "",
		DB:       0,
	})

	re := regexp.MustCompile(cookie_prefix + `.*`)

	t, err := tail.TailFile(wp_log, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		//log.Println(line.Text)
		line_slice := strings.Split(line.Text, ";")
		for _, value := range line_slice {
			//log.Println(value)
			if strings.Contains(value, cookie_prefix) {
				cookie_log = strings.TrimSpace(value)
				//log.Println("Success, cookie found: " + cookie_log)
				cookie := strings.TrimPrefix(strings.Join(re.FindAllString(cookie_log, -1), ""), cookie_prefix)
				//log.Println(cookie)
				uuid := uuid.New().String()
				key := "wp-cookie-" + uuid
				log.Println("cookie found, setting " + key + " key in redis...")
				set, err := rdb.SetNX(ctx, key, cookie, redis_ttl_time).Result()
				if err != nil {
					log.Fatal(err)
				}
				log.Println(set)
				break
			}
		}
		err := os.Truncate(wp_log, 0)
		if err != nil {
			log.Fatal(err)

		}
	}
}
