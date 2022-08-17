package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/apcera/termtables"
	MoonPhase "github.com/janczer/goMoonPhase"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kyokomi/emoji"
	"github.com/mshafiee/swephgo"
)

type Moon struct {
	date      time.Time
	unix      int64
	phase     float64
	phaseName string
	emoji     string
}

var (
	loc         string
	city        = os.Getenv("CITY")
	houseSystem = os.Getenv("HOUSE_SYSTEM")
	location    *time.Location
	swisspath   = os.Getenv("SWISSPATH")
)

func init() {
	location, _ = time.LoadLocation(city)
	swephgo.SetEphePath([]byte(swisspath))
}

func main() {
	// p := message.NewPrinter(language.Polish)

	now := time.Now()
	start := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	end := now.AddDate(1, 0, 1) // look ahead up to 1 year and 1 day
	table1 := termtables.CreateTable()
	table1.AddHeaders("New Moon", "", "Full Moon", "")

	for d := start; d.After(end) == false; {
		newMoon, fullMoon := moonPhase(d)
		d = fullMoon.date.AddDate(0, 0, 13)
		// mph := MoonPhase.New(newMoon.date)
		// emoji.Printf("%s %s %s\n", newMoon.date.Format(time.RFC822), moonEmoji(phaseName), phaseName)
		// s := p.Sprintf("Age: %.1f Ilumination: %.3f Phase: %.3f Longitude: %.3f Distance: %.0f km\n", mph.Age(), mph.Illumination(), mph.Phase(), mph.Longitude(), mph.Distance())
		// fmt.Println(s)

		// mph = MoonPhase.New(fullMoon.date)
		// phaseName = mph.PhaseName()
		// emoji.Printf("%s %s %s\n", fullMoon.date.Format(time.RFC822), moonEmoji(phaseName), phaseName)
		// s = p.Sprintf("Age: %.1f Ilumination: %.3f Phase: %.3f Longitude: %.3f Distance: %.0f km\n", mph.Age(), mph.Illumination(), mph.Phase(), mph.Longitude(), mph.Distance())
		// fmt.Println(s)
		table1.AddRow(newMoon.date.Format(time.RFC822), newMoon.emoji, fullMoon.date.Format(time.RFC822), fullMoon.emoji)
	}
	fmt.Println(table1.Render())
}

/*
	1) Find nearest New Moon within 29 days from now
	2) Find exact time of new moon (up to a minute)
	3) Jump 14 days
	4) Find exact time of Full Moon (up to a minute)
	5) Jump 14 days
	6) Repeat 2
*/
func moonPhase(start time.Time) (*Moon, *Moon) {
	end := start.AddDate(0, 1, 1) // look ahead up to 2 years and 1 day
	newMoon := Moon{
		date:  start,
		unix:  start.Unix(),
		phase: 1.0,
	}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		phase, _ := Phase(d, swephgo.SeMoon)
		// s := p.Sprintf("Phase: %.3f Prev: %.3f", phase, newMoon.phase)
		// fmt.Println(s)
		if math.Abs(1.0-phase) < math.Abs(1.0-newMoon.phase) {
			break
		}
		newMoon.date = d
		newMoon.unix = d.Unix()
		newMoon.phase = phase
	}
	// s := p.Sprintf("New Moon aprox. %s", newMoon.date.Format(time.RFC3339))
	// fmt.Println(s)
	start = newMoon.date.Add(time.Hour * -24)
	end = newMoon.date.Add(time.Hour * 24)
	newMoon = Moon{
		date:  start,
		unix:  start.Unix(),
		phase: 1.0,
	}
	for d := start; d.After(end) == false; d = d.Add(time.Minute * 1) {
		phase, _ := Phase(d, swephgo.SeMoon)
		// s := p.Sprintf("Phase: %.10f Prev: %.10f", phase, newMoon.phase)
		// fmt.Println(s)
		if math.Abs(1.0-phase) < math.Abs(1.0-newMoon.phase) {
			break
		}
		newMoon.date = d
		newMoon.unix = d.Unix()
		newMoon.phase = phase
	}
	// s := p.Sprintf("New Moon exact. %s", newMoon.date.Format(time.RFC3339))
	// fmt.Println(s)
	mph := MoonPhase.New(newMoon.date)
	newMoon.phaseName = mph.PhaseName()
	newMoon.emoji = emoji.Sprintf("%s", moonEmoji(newMoon.phaseName))

	start = newMoon.date.AddDate(0, 0, 13)
	end = newMoon.date.AddDate(0, 0, 15) // look ahead up to 2 days
	fullMoon := Moon{
		date:  start,
		unix:  start.Unix(),
		phase: 0.0,
	}
	for d := start; d.After(end) == false; d = d.Add(time.Minute * 1) {
		phase, _ := Phase(d, swephgo.SeMoon)
		// s := p.Sprintf("Phase: %.10f Prev: %.10f", phase, fullMoon.phase)
		// fmt.Println(s)
		if math.Abs(phase) < math.Abs(fullMoon.phase) {
			break
		}
		fullMoon.date = d
		fullMoon.unix = d.Unix()
		fullMoon.phase = phase
	}
	// s = p.Sprintf("Full Moon exact. %s", fullMoon.date.Format(time.RFC3339))
	// fmt.Println(s)
	mph = MoonPhase.New(fullMoon.date)
	fullMoon.phaseName = mph.PhaseName()
	fullMoon.emoji = emoji.Sprintf("%s", moonEmoji(fullMoon.phaseName))

	return &newMoon, &fullMoon
}

func moonEmoji(icon string) string {
	switch icon {
	case "New Moon":
		return ":new_moon:"
	case "Waxing Crescent":
		return ":waxing_crescent_moon:"
	case "First Quarter":
		return ":first_quarter_moon:"
	case "Waxing Gibbous":
		return ":waxing_gibbous_moon:"
	case "Full Moon":
		return ":full_moon:"
	case "Waning Gibbous":
		return ":waning_gibbous_moon:"
	case "Third Quarter":
		return ":last_quarter_moon:"
	case "Waning Crescent":
		return ":waning_crescent_moon:"
	default:
		return ":star:"
	}
}
