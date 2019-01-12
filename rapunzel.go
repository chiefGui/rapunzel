package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/fatih/structs"
)

const outputMobs = "./output/mob_db.conf"
const inputMobs = "./input/mob_db.txt"

func main() {
	content := openMobDb()
	parsedMobs := parseMobs(content)

	err := saveMobsOnDisk(parsedMobs)
	checkErr(err)
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func openMobDb() string {
	content, err := ioutil.ReadFile(inputMobs)
	checkErr(err)

	return string(content)
}

func parseMobs(mobs string) []Mob {
	allMobs := strings.Split(mobs, "\n")
	allMobs = allMobs[:len(allMobs)-1]
	parsedMobs := []Mob{}

	for _, mob := range allMobs {
		parsedMob, err := parseMob(mob)

		if err != nil {
			parsedMobs = append(parsedMobs, parsedMob)
		}
	}

	return parsedMobs
}

func parseMob(mob string) (Mob, error) {
	mobProperties := strings.Split(mob, ",")
	parsedMob := Mob{}

	if strings.Contains(mobProperties[0], "//") {
		color.Blue("- Skipped one comment.")
		return parsedMob, errors.New("mob is commented")
	}

	parsedMob.ID = mobProperties[0]
	parsedMob.SpriteName = mobProperties[1]
	parsedMob.Name = mobProperties[2]
	parsedMob.JName = mobProperties[2]

	color.Green("\"%s\": converted successfully!\n", parsedMob.Name)

	return parsedMob, nil
}

func tryToDeleteMobsOnDisk() bool {
	lock.Lock()
	defer lock.Unlock()

	err := os.Remove(outputMobs)
	if err != nil {
		return true
	}

	return false
}

func saveMobsOnDisk(parsedMobs []Mob) error {
	tryToDeleteMobsOnDisk()

	lock.Lock()
	defer lock.Unlock()

	file, err := os.Create(outputMobs)
	checkErr(err)

	defer file.Close()

	mobsInHerculesSyntax := convertParsedMobsToHerculesSyntax(parsedMobs)

	_, err = file.WriteString(mobsInHerculesSyntax)

	fmt.Println("")
	color.Cyan("mob_db.conf generated successfully. Please check your 'output/' folder.")

	return err
}

func convertParsedMobsToHerculesSyntax(parsedMobs []Mob) string {
	mobsInHerculesSyntax := []string{}

	for index, mob := range parsedMobs {
		isLast := false

		if index == len(parsedMobs)-1 {
			isLast = true
		}

		props := []string{}
		keysNames := structs.Names(&Mob{})

		for _, keyName := range keysNames {
			props = append(props, propObject(keyName, getValueFromKey(mob, keyName)))
		}

		shape := append(startObject(), strings.Join(props, ""), endObject(isLast))

		strings.Join(shape, "")

		mobsInHerculesSyntax = append(mobsInHerculesSyntax, strings.Join(shape, ""))
	}

	return strings.Join(mobsInHerculesSyntax, "")
}

func startObject() []string {
	return []string{
		"{ \n",
	}
}

func propObject(key string, value string) string {
	prop := []string{
		"	",
		key,
		": ",
		value,
		"\n",
	}

	return strings.Join(prop, "")
}

func endObject(isLast bool) string {
	s := []string{}

	if isLast {
		s = []string{
			"}",
		}
	} else {
		s = []string{
			"},",
		}
	}

	return strings.Join(s, "")
}

func getValueFromKey(v interface{}, keyName string) string {
	r := reflect.ValueOf(v)
	value := reflect.Indirect(r).FieldByName(keyName)

	return value.String()
}

var lock sync.Mutex

// Mob is a parsed mob.
type Mob struct {
	ID          string
	SpriteName  string
	Name        string
	JName       string
	Lv          string
	Hp          string
	Sp          string
	Exp         string
	JExp        string
	AttackRange string
	Attack      string
	Def         string
	Mdef        string
	Str         string
	Agi         string
	Vit         string
	Int         string
	Dex         string
	Luk         string
	Range2      string
}
