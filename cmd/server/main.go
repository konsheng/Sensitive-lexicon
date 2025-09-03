package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"sensitive-lexicon/internal/detect"
	"sensitive-lexicon/internal/lexicon"
	"strconv"
	"time"
)

type User struct {
	Name       string    `json:"name" validate:"min=5,max=20"`
	Age        int       `json:"age" validate:"gte=18"`
	Enrollment time.Time `json:"enrollment" validate:"before_today"`
	Graduation time.Time `json:"graduation" validate:"gtfield=Enrollment"`
}

// BeforeToday 验证日期是否在今天之前
func BeforeToday(fl validator.FieldLevel) bool {
	fieldTime, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return fieldTime.Before(time.Now())
}
func main() {
	lexiconDir := getenv("LEXICON_DIR", "Vocabulary")
	minNgram := getenvInt("FUZZY_MIN_NGRAM", 2)
	maxNgram := getenvInt("FUZZY_MAX_NGRAM", 10)
	maxDistance := getenvInt("FUZZY_MAX_DISTANCE", 1)

	store := lexicon.NewStore()
	if err := store.LoadFromDir(lexiconDir); err != nil {
		log.Fatalf("failed to load lexicon: %v", err)
	}

	service := detect.NewService(store)
	service.SetFuzzyConfig(detect.FuzzyConfig{MinNgramLen: minNgram, MaxNgramLen: maxNgram, MaxDistance: maxDistance})

	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/detect", func(c *fiber.Ctx) error {
		var req detect.DetectRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		res := service.Detect(req)
		return c.JSON(res)
	})

	app.Post("/contains", func(c *fiber.Ctx) error {
		var req detect.ContainsRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		res := service.Contains(req)
		return c.JSON(res)
	})

	app.Post("/reload", func(c *fiber.Ctx) error {
		if err := store.LoadFromDir(lexiconDir); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		stats := store.Stats()
		return c.JSON(stats)
	})

	port := getenv("PORT", "8080")
	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
