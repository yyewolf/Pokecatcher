package main

import "strconv"

const (
	filterLevel    = iota //0
	filterHP              //1
	filterAttack          //2
	filterDefense         //3
	filterSpAttack        //4
	filterSpDef           //5
	filterSpeed           //6
	filterIV              //7
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
		break
	case filterEqual:
		return n1 == n2
		break
	case filterAbove:
		return n1 > n2
		break
	}
	return false
}

func filterDo(t, c int, n float64, infos PokeInfoParsed) bool {
	switch t {
	case filterLevel:
		lvl, _ := strconv.Atoi(infos.Level)
		return filterCheck(c, float64(lvl), n)
		break
	case filterHP:
		return filterCheck(c, float64(infos.IVs[0]), n)
		break
	case filterAttack:
		return filterCheck(c, float64(infos.IVs[1]), n)
		break
	case filterDefense:
		return filterCheck(c, float64(infos.IVs[2]), n)
		break
	case filterSpAttack:
		return filterCheck(c, float64(infos.IVs[3]), n)
		break
	case filterSpDef:
		return filterCheck(c, float64(infos.IVs[4]), n)
		break
	case filterSpeed:
		return filterCheck(c, float64(infos.IVs[5]), n)
		break
	case filterIV:
		return filterCheck(c, infos.TotalIV, n)
		break
	}
	return false
}

func allFilters(infos PokeInfoParsed) bool {
	r1 := false                          // Final Result
	for i := range Config.EveryFilters { //This is the OR part
		r2 := true                                      // Middle Result
		for j := range Config.EveryFilters[i].Filters { //This is the AND part
			t := Config.EveryFilters[i].Filters[j].ToCheck
			c := Config.EveryFilters[i].Filters[j].Operation
			n := Config.EveryFilters[i].Filters[j].ComparedTo
			if !filterDo(t, c, n, infos) {
				r2 = false
			}
		}
		r1 = r2
	}
	return r1
}
