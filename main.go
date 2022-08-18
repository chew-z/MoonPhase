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
	age         float64
	date        time.Time
	unix        int64
	phase       float64
	phaseName   string
	ilumination float64
	emoji       string
	zodiac      string
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

	now := time.Now()
	start := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, location)
	end := now.AddDate(1, 0, 1) // look ahead up to 1 year and 1 day
	table1 := termtables.CreateTable()
	table1.AddHeaders("New Moon", "", "Full Moon", "")

	for d := start; d.After(end) == false; {
		newMoon, fullMoon := moonPhase(d)
		table1.AddRow(newMoon.date.Format(time.RFC822), newMoon.zodiac, fullMoon.date.Format(time.RFC822), fullMoon.zodiac)
		d = fullMoon.date.AddDate(0, 0, 14)
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
	end := start.AddDate(0, 0, 31) // look ahead up to 1 month and 1 day
	newMoon := Moon{
		date:  start,
		unix:  start.Unix(),
		phase: 1.0,
	}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		phase, _ := Phase(d, swephgo.SeMoon)
		if math.Abs(1.0-phase) < math.Abs(1.0-newMoon.phase) {
			break
		}
		newMoon.date = d
		newMoon.unix = d.Unix()
		newMoon.phase = phase
	}
	start = newMoon.date.Add(time.Hour * -24)
	end = newMoon.date.Add(time.Hour * 24)
	ps, _ := Phase(start, swephgo.SeMoon)
	startMoon := Moon{
		date:  start,
		unix:  start.Unix(),
		phase: ps,
	}
	pe, _ := Phase(end, swephgo.SeMoon)
	endMoon := Moon{
		date:  end,
		unix:  end.Unix(),
		phase: pe,
	}
	newMoon = binarySearch(startMoon, endMoon, false)
	mph := MoonPhase.New(newMoon.date)
	newMoon.phaseName = mph.PhaseName()
	newMoon.emoji = emoji.Sprintf("%s", moonEmoji(newMoon.phaseName))
	newMoon.ilumination = mph.Illumination()
	newMoon.age = mph.Age()
	newMoon.zodiac = mph.ZodiacSign()

	start = newMoon.date.AddDate(0, 0, 14)
	end = newMoon.date.AddDate(0, 0, 16) // look ahead up to 2 days
	ps, _ = Phase(start, swephgo.SeMoon)
	startMoon = Moon{
		date:  start,
		unix:  start.Unix(),
		phase: ps,
	}
	pe, _ = Phase(end, swephgo.SeMoon)
	endMoon = Moon{
		date:  end,
		unix:  end.Unix(),
		phase: pe,
	}
	fullMoon := binarySearch(startMoon, endMoon, true)
	mph = MoonPhase.New(fullMoon.date)
	fullMoon.phaseName = mph.PhaseName()
	fullMoon.emoji = emoji.Sprintf("%s", moonEmoji(fullMoon.phaseName))
	fullMoon.ilumination = mph.Illumination()
	fullMoon.age = mph.Age()
	fullMoon.zodiac = mph.ZodiacSign()

	return &newMoon, &fullMoon
}

func binarySearch(start Moon, end Moon, fullMoon bool) Moon {
	half := end.date.Sub(start.date).Seconds() / 2
	mDate := start.date.Add(time.Second * time.Duration(half))
	phase, _ := Phase(mDate, swephgo.SeMoon)

	// p := message.NewPrinter(language.Polish)
	// sp := p.Sprintf("Start: %s Mid: %s End: %s", start.date.Format(time.RFC822), mDate.Format(time.RFC822), end.date.Format(time.RFC822))
	// fmt.Println(sp)

	newStart := start
	newEnd := end
	middle := Moon{
		date:  mDate,
		unix:  mDate.Unix(),
		phase: phase,
	}
	// sp = p.Sprintf("Start: %.15f Phase: %.15f End: %.15f", newStart.phase, middle.phase, newEnd.phase)
	// fmt.Println(sp)

	if fullMoon {
		if start.phase < end.phase {
			newStart = middle
		} else {
			newEnd = middle
		}
	} else {
		if end.phase < start.phase {
			newStart = middle
		} else {
			newEnd = middle
		}
	}
	if newEnd.date.Sub(newStart.date).Minutes() < 1.0 {
		return newEnd
	}
	return binarySearch(newStart, newEnd, fullMoon)
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
