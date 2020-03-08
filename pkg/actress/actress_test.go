package actress

import (
	"fmt"
	"testing"

	"github.com/smartystreets/assertions/should"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ylqjgm/AVMeta/pkg/util"
)

var (
	act   *Actress
	first = true
)

func TestNewActress(t *testing.T) {
	Convey("New Actress", t, func() {
		monkey.Patch(util.GetConfig, func() (*util.ConfigStruct, error) {
			return &util.ConfigStruct{}, fmt.Errorf("error")
		})
		defer monkey.UnpatchAll()
		act = NewActress()
		So(act, ShouldBeNil)
	})
}

func TestActress_Fetch(t *testing.T) {
	Convey("Actress Fetch", t, func() {
		monkey.Patch(util.GetConfig, func() (*util.ConfigStruct, error) {
			return &util.ConfigStruct{
				Base: util.BaseStruct{
					Proxy: "",
				},
				Site: util.SiteStruct{
					JavBus: "https://www.javbus.com/",
					JavDB:  "https://javdb4.com/",
				},
			}, nil
		})
		monkey.Patch(util.SavePhoto, func(_, _, _ string, _ bool) error {
			return nil
		})
		defer monkey.UnpatchAll()

		act = NewActress()

		Convey("JavBUS", func() {
			monkey.Patch(util.GetConfig, func() (*util.ConfigStruct, error) {
				return &util.ConfigStruct{
					Base: util.BaseStruct{
						Proxy: "",
					},
					Site: util.SiteStruct{
						JavBus: "https://www.javbus.com/",
						JavDB:  "https://javdb4.com/",
					},
				}, nil
			})
			monkey.Patch(JavBUS, func(_, _ string, _ int, _ bool) (_ map[string]string, _ bool, _ error) {
				actors := make(map[string]string)
				actors["北条麻妃"] = "https://us.netcdn.space/tokyohot/media/cast/2395/thumbnail.jpg"

				return actors, false, nil
			})
			defer monkey.UnpatchAll()

			err := act.Fetch(JAVBUS, 1, true)
			So(err, ShouldBeNil)
		})

		Convey("JavDB", func() {
			monkey.Patch(JavDB, func(_, _ string, _ int, _ bool) (_ map[string]string, _ bool, _ error) {
				actors := make(map[string]string)
				actors["北条麻妃"] = "https://us.netcdn.space/tokyohot/media/cast/2395/thumbnail.jpg"

				return actors, false, nil
			})
			defer monkey.UnpatchAll()

			err := act.Fetch(JAVDB, 1, true)
			So(err, ShouldBeNil)
		})

		Convey("Site error", func() {
			err := act.Fetch("AAA", 1, true)
			So(err, ShouldBeError)
		})

		Convey("Fetch error", func() {
			monkey.Patch(JavDB, func(_, _ string, _ int, _ bool) (_ map[string]string, _ bool, _ error) {
				return nil, false, fmt.Errorf("error")
			})
			defer monkey.UnpatchAll()

			err := act.Fetch(JAVDB, 1, true)
			So(err, ShouldBeError)
		})

		Convey("Next page", func() {
			monkey.Patch(JavDB, func(_, _ string, _ int, _ bool) (_ map[string]string, _ bool, _ error) {
				if first {
					first = false
					return nil, true, nil
				}

				return nil, false, fmt.Errorf("test finish")
			})
			defer monkey.UnpatchAll()

			err := act.Fetch(JAVDB, 1, true)
			So(err, should.BeError, "test finish")
		})
	})
}

func TestActress_Put(t *testing.T) {
	Convey("Put actress", t, func() {
	})
}
