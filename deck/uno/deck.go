package uno

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/KirCute/sim-board"
)

type Uno struct {
	rest     []sim_board.Card
	shuffled bool
	*Params
}

func (p *Uno) Type() string {
	return Name
}

func (p *Uno) Name() string {
	if len(p.CustomName) == 0 {
		return Name
	}
	return p.CustomName
}

func (p *Uno) RestLen() int {
	return len(p.rest)
}

func (p *Uno) MaxLen() int {
	cntPerColor := p.CountRank + p.CountColoredReverse + p.CountColoredSkip + p.CountColoredApp2 + p.CountColoredApp4 +
		p.CountColoredApp6 + p.CountColoredApp8 + p.CountColoredApp10 + p.CountColoredDiscardAll + p.CountColoredSwap +
		p.CountColoredSkipAll + p.CountColoredBlank
	cntBlack := p.CountTrans + p.CountBlackReverse + p.CountBlackSkip + p.CountBlackApp2 + p.CountBlackApp4 +
		p.CountBlackApp6 + p.CountBlackApp8 + p.CountBlackApp10 + p.CountBlackDiscardAll + p.CountBlackSwap +
		p.CountBlackSkipAll + p.CountBlackBlank
	return cntPerColor*p.CountColor*p.Count + cntBlack
}

func (p *Uno) Return(card sim_board.Card) {
	p.rest = append(p.rest, card)
	p.shuffled = false
}

func (p *Uno) Draw(count int) []sim_board.Card {
	if !p.shuffled {
		rand.Shuffle(len(p.rest), func(i, j int) {
			p.rest[i], p.rest[j] = p.rest[j], p.rest[i]
		})
		p.shuffled = true
	}
	ret := p.rest[:count]
	p.rest = p.rest[count:]
	return ret
}

func GetHTML(card sim_board.Card) (string, bool) {
	color, content, ok := strings.Cut(string(card), "-")
	if !ok {
		return "", false
	}
	var cardContent string
	if text, ok := strings.CutPrefix(content, "plain-"); ok {
		cardContent = fmt.Sprintf(`<span style="position: absolute; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%); color: %s; font-size: 30px; font-weight: bold; font-family: Arial, sans-serif">%s</span>`, color, text)
	} else {
		switch content {
		case "skip":
			cardContent = fmt.Sprintf(`
<svg style="position: absolute; width: 32px; height: 32px; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%)" viewBox="0 0 100 100">
  <circle cx="50" cy="50" r="40" fill="none" stroke="%s" stroke-width="15" />
  <line x1="78" y1="22" x2="22" y2="78" stroke="%s" stroke-width="15" stroke-linecap="round" />
</svg>
`, color, color)
		case "reverse":
			cardContent = fmt.Sprintf(`
<svg style="position: absolute; width: 32px; height: 32px; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%)" viewBox="0 0 100 100">
  <polygon points="29,68 26,64 24,59 23,54 24,49 26,44 29,40 56,13 48,5 78,5 78,35 70,27" fill="%s" />
  <polygon points="70,33 73,37 75,42 76,47 75,52 73,57 70,61 43,88 51,96 21,96 21,66 29,74" fill="%s" />
</svg>
`, color, color)
		case "skipall":
			cardContent = fmt.Sprintf(`
<svg style="position: absolute; width: 32px; height: 32px; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%)" viewBox="0 0 170 170">
  <circle cx="50" cy="120" r="40" fill="none" stroke="%s" stroke-width="15" />
  <line x1="78" y1="92" x2="22" y2="148" stroke="%s" stroke-width="15" stroke-linecap="round" />
  <circle cx="120" cy="50" r="40" fill="none" stroke="%s" stroke-width="15" />
  <line x1="148" y1="22" x2="92" y2="78" stroke="%s" stroke-width="15" stroke-linecap="round" />
</svg>
`, color, color, color, color)
		case "discardall":
			cardContent = fmt.Sprintf(`
<svg style="position: absolute; width: 32px; height: 32px; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%)" viewBox="0 0 95 95">
  <rect x="17" y="32" width="41" height="51" rx="4" ry="4" fill="%s" stroke="none"/>
  <rect x="12" y="27" width="51" height="61" rx="7" ry="7" fill="none" stroke="white" stroke-width="5"/>
  <rect x="10" y="25" width="55" height="65" rx="7" ry="7" fill="none" stroke="%s" stroke-width="2"/>
  <polygon points="65,25 80,25 75,20 85,10 80,5 70,15 65,10" fill="%s" />
</svg>
`, color, color, color)
		case "swap":
			cardContent = fmt.Sprintf(`
<svg style="position: absolute; width: 32px; height: 32px; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%)" viewBox="0 0 100 100">
  <polygon points="29,68 26,64 24,59 23,54 24,49 26,44 29,40 56,13 48,5 78,5 78,35 70,27" fill="%s" />
  <polygon points="70,33 73,37 75,42 76,47 75,52 73,57 70,61 43,88 51,96 21,96 21,66 29,74" fill="%s" />
  <rect x="34" y="29" width="32" height="42" rx="4" ry="4" fill="%s" stroke="none"/>
  <rect x="32" y="27" width="36" height="46" rx="7" ry="7" fill="none" stroke="white" stroke-width="5"/>
  <rect x="30" y="25" width="40" height="50" rx="7" ry="7" fill="none" stroke="%s" stroke-width="2"/>
</svg>
`, color, color, color, color)
		case "blank":
			cardContent = ""
		case "trans":
			cardContent = fmt.Sprintf(`
<div style="position: absolute; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%) rotate(-45deg)">
  <div style="position: absolute; left: 50%%; top: 50%%; width: 38px; height: 38px; background: conic-gradient(%s 0deg 90deg, %s 90deg 180deg, %s 180deg 270deg, %s 270deg 360deg); border-radius: 50%%; transform: translate(-50%%, -50%%) scaleX(1.6585)"></div>
</div>
`, COLORS[0], COLORS[3], COLORS[2], COLORS[1])
		default:
			cardContent = fmt.Sprintf(`<div style="position: absolute; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%); color: %s">%s</div>`, color, content)
		}
	}
	return fmt.Sprintf(`
<div style="position: relative; width: 65px; aspect-ratio: 0.7222; box-sizing: border-box; border: 5px solid white; background-color: %s; border-radius: 10px; overflow: hidden">
  <div style="position: absolute; width: 68px; height: 41px; background-color: white; border-radius: 50%%; left: 50%%; top: 50%%; transform: translate(-50%%, -50%%) rotate(-45deg)"></div>
  %s
</div>
`, color, cardContent), true
}
