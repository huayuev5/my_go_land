package strategy

import (
	"fmt"
	"testing"

	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/g"
	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/scheme"
)

func TestGetMyStrategy(t *testing.T) {
	fmt.Println("Now Test GetLocalStrategy:")
	data, err := getFileStrategy()
	if err == nil {
		fmt.Println("Result:")
		for _, x := range data {
			fmt.Printf("    %v\n", x)
		}
	} else {
		fmt.Println("Something Error:")
		fmt.Println(err)
	}
}

func TestPatternParse(t *testing.T) {
	fmt.Println("Now Test PatternParse:")
	var a scheme.Strategy
	a.Pattern = "memeda"
	parsePattern([]*scheme.Strategy{&a})
	fmt.Printf("a.pat:[%s], a.ex[%s]\n", a.Pattern, a.Exclude)
	a.Pattern = "```EXCLUDE```memeda"
	parsePattern([]*scheme.Strategy{&a})
	fmt.Printf("a.pat:[%s], a.ex[%s]\n", a.Pattern, a.Exclude)
	a.Pattern = "memeda```EXCLUDE```"
	parsePattern([]*scheme.Strategy{&a})
	fmt.Printf("a.pat:[%s], a.ex[%s]\n", a.Pattern, a.Exclude)
}

func TestGetFromFile(t *testing.T) {
	fmt.Println("Now Test Get Strategy From File:")
	g.InitStrategyFile()
	sts, err := getFileStrategy()
	if err != nil {
		fmt.Println("Read Error: %v", err)
	}
	for _, st := range sts {
		fmt.Println(st)
	}
}
