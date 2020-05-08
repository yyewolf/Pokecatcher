package main

import "strconv"

const (
	filterLevel       = iota //0
	filterHP                 //1
	filterAttack             //2
	filterDefense            //3
	filterSpAttack           //4
	filterSpDef              //5
	filterSpeed              //6
	filterIV                 //7
	filterIsAlolan           //8
	filterIsGalarian         //9
	filterIsLegendary        //10
	filterIsShiny            //11
)

const (
	filterUnder = iota //0
	filterEqual        //1
	filterAbove        //2
)

//FilterStruct = Save filter into the config file
type filterStruct struct {
	ToCheck    int     `json:"Checking"`
	Operation  int     `json:"Operation"`
	ComparedTo float64 `json:"Value"`
}

//FilterStruct = Save filter into the config file
type customFilterStruct struct {
	Filters []filterStruct `json:"Conditions"`
}

func filterCheck(t int, n1, n2 float64) bool {
	switch t {
	case filterUnder:
		return n1 < n2
	case filterEqual:
		return n1 == n2
	case filterAbove:
		return n1 > n2
	}
	return false
}

func filterDo(t, c int, n float64, infos pokeInfoParsed) bool {
	switch t {
	case filterLevel:
		lvl, _ := strconv.Atoi(infos.Level)
		return filterCheck(c, float64(lvl), n)
	case filterHP:
		return filterCheck(c, float64(infos.IVs[0]), n)
	case filterAttack:
		return filterCheck(c, float64(infos.IVs[1]), n)
	case filterDefense:
		return filterCheck(c, float64(infos.IVs[2]), n)
	case filterSpAttack:
		return filterCheck(c, float64(infos.IVs[3]), n)
	case filterSpDef:
		return filterCheck(c, float64(infos.IVs[4]), n)
	case filterSpeed:
		return filterCheck(c, float64(infos.IVs[5]), n)
	case filterIV:
		return filterCheck(c, infos.TotalIV, n)
	case filterIsAlolan:
		return infos.isAlolan
	case filterIsGalarian:
		return infos.isGalarian
	case filterIsLegendary:
		//Searches in the name list
		ok := false
		for i := range legendaries {
			if infos.Name == legendaries[i] {
				return !ok
			}
		}
		return ok
	case filterIsShiny:
		return infos.isShiny
	}
	return false
}

func allFilters(infos pokeInfoParsed) bool {
	r1 := false                          // Final Result
	for i := range config.EveryFilters { //This is the OR part
		r2 := true                                      // Middle Result
		for j := range config.EveryFilters[i].Filters { //This is the AND part
			t := config.EveryFilters[i].Filters[j].ToCheck
			c := config.EveryFilters[i].Filters[j].Operation
			n := config.EveryFilters[i].Filters[j].ComparedTo
			if !filterDo(t, c, n, infos) {
				r2 = false
			}
		}
		r1 = r2
	}
	return r1
}
