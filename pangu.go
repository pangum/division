package division

import `github.com/storezhang/pangu`

func init() {
	if err := pangu.New().Provides(newDivision); nil != err {
		panic(err)
	}
}
