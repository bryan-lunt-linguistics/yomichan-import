/*
 * Copyright (c) 2016 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"regexp"
	"strings"
)

type daijirinExtractor struct {
	partsExp   *regexp.Regexp
	phonExp    *regexp.Regexp
	variantExp *regexp.Regexp
	annotExp   *regexp.Regexp
	v5Exp      *regexp.Regexp
	v1Exp      *regexp.Regexp
}

func makeDaijirinExtractor() epwingExtractor {
	return &daijirinExtractor{
		partsExp:   regexp.MustCompile(`([^（【〖]+)(?:【(.*)】)?(?:〖(.*)〗)?(?:（(.*)）)?`),
		phonExp:    regexp.MustCompile(`[-・]+`),
		variantExp: regexp.MustCompile(`\((.*)\)`),
		annotExp:   regexp.MustCompile(`（(.*)）`),
		v5Exp:      regexp.MustCompile(`(動.五)|(動..二)`),
		v1Exp:      regexp.MustCompile(`動..一`),
	}
}

func (e *daijirinExtractor) extractTerms(entry epwingEntry) []dbTerm {
	matches := e.partsExp.FindStringSubmatch(entry.Heading)
	if matches == nil {
		return nil
	}

	var expressions, readings []string
	if expression := matches[2]; len(expression) > 0 {
		expression = e.annotExp.ReplaceAllLiteralString(expression, "")
		for _, split := range strings.Split(expression, "・") {
			splitInc := e.variantExp.ReplaceAllString(split, "$1")
			expressions = append(expressions, splitInc)
			if split != splitInc {
				splitExc := e.variantExp.ReplaceAllLiteralString(split, "")
				expressions = append(expressions, splitExc)
			}
		}
	}

	if reading := matches[1]; len(reading) > 0 {
		reading = e.phonExp.ReplaceAllLiteralString(reading, "")
		readings = append(readings, reading)
	}

	var tags []string
	for _, split := range strings.Split(entry.Text, "\n") {
		if matches := e.annotExp.FindStringSubmatch(split); matches != nil {
			for _, tag := range strings.Split(matches[1], "・") {
				tags = append(tags, tag)
			}
		}
	}

	var terms []dbTerm
	if len(expressions) == 0 {
		for _, reading := range readings {
			term := dbTerm{
				Expression: reading,
				Glossary:   []string{entry.Text},
			}

			e.exportTags(&term, tags)
			e.exportRules(&term, tags)

			terms = append(terms, term)
		}

	} else {
		for _, expression := range expressions {
			for _, reading := range readings {
				term := dbTerm{
					Expression: expression,
					Reading:    reading,
					Glossary:   []string{entry.Text},
				}

				e.exportTags(&term, tags)
				e.exportRules(&term, tags)

				terms = append(terms, term)
			}
		}
	}

	return terms
}

func (*daijirinExtractor) extractKanji(entry epwingEntry) []dbKanji {
	return nil
}

func (e *daijirinExtractor) exportRules(term *dbTerm, tags []string) {
	for _, tag := range tags {
		if tag == "形" {
			term.addTags("adj-i")
			term.addRules("adj-i")
		} else if e.v5Exp.MatchString(tag) {
			term.addTags("v5")
			term.addRules("v5")
		} else if e.v1Exp.MatchString(tag) {
			term.addTags("v1")
			term.addRules("v1")
		}
	}
}

func (*daijirinExtractor) getRevision() string {
	return "daijirin:1"
}

func (e *daijirinExtractor) exportTags(term *dbTerm, tags []string) {
	parsed := []string{
		"並立助",
		"代",
		"係助",
		"副",
		"副助",
		"助動",
		"動",
		"動ア上一",
		"動ア下一",
		"動ア下二",
		"動ア五［ハ四］",
		"動カ上一",
		"動カ上二",
		"動カ下一",
		"動カ下二",
		"動カ五",
		"動カ五［四］",
		"動カ四",
		"動カ変",
		"動ガ上一",
		"動ガ上二",
		"動ガ下一",
		"動ガ下二",
		"動ガ五",
		"動ガ五［四］",
		"動ガ四",
		"動サ上一",
		"動サ下一",
		"動サ下二",
		"動サ五",
		"動サ五［四］",
		"動サ四",
		"動サ変",
		"動サ特活",
		"動ザ上一",
		"動ザ上二",
		"動ザ下一",
		"動ザ下二",
		"動タ上一",
		"動タ上二",
		"動タ下一",
		"動タ下二",
		"動タ五［四］",
		"動タ四",
		"動ダ上二",
		"動ダ下一",
		"動ダ下二",
		"動ナ上一",
		"動ナ下一",
		"動ナ下二",
		"動ナ五",
		"動ナ五［四］",
		"動ハ上一",
		"動ハ上二",
		"動ハ下一",
		"動ハ下二",
		"動ハ四",
		"動ハ特活",
		"動バ上一",
		"動バ上二",
		"動バ下一",
		"動バ下二",
		"動バ五［四］",
		"動バ四",
		"動マ上一",
		"動マ上二",
		"動マ下一",
		"動マ下二",
		"動マ五",
		"動マ五［四］",
		"動マ四",
		"動マ特活",
		"動マ特活",
		"動ヤ上一",
		"動ヤ上二",
		"動ヤ下二",
		"動ラ上一",
		"動ラ上二",
		"動ラ下一",
		"動ラ下二",
		"動ラ五",
		"動ラ五［四］",
		"動ラ四",
		"動ラ変",
		"動ラ特活",
		"動ワ上一",
		"動ワ上二",
		"動ワ下一",
		"動ワ下二",
		"動ワ五",
		"動ワ五［ハ四］",
		"動五［四］",
		"動特活",
		"動詞五［四］段型活用",
		"名",
		"形",
		"形ク",
		"形シク",
		"形動",
		"形動タリ",
		"形動ナリ",
		"感",
		"接助",
		"接尾",
		"接続",
		"接頭",
		"枕詞",
		"格助",
		"終助",
		"連体",
		"連語",
		"間投助",
	}

	for _, tag := range tags {
		for _, p := range parsed {
			if tag == p {
				term.addTags(tag)
			}
		}
	}
}

func (*daijirinExtractor) getFontNarrow() map[int]string {
	return map[int]string{
		49441: "á",
		49442: "à",
		49443: "â",
		49444: "ä",
		49445: "ã",
		49446: "ā",
		49447: "é",
		49448: "è",
		49449: "ê",
		49450: "ë",
		49451: "ē",
		49452: "í",
		49453: "î",
		49454: "ï",
		49455: "ñ",
		49456: "ó",
		49457: "ò",
		49458: "ô",
		49459: "ö",
		49460: "ř",
		49461: "ú",
		49462: "ü",
		49463: "~",
		49464: "ç",
		49465: "ˇ",
		49466: "ɡ",
		49467: "ŋ",
		49468: "ʒ",
		49469: "ʃ",
		49470: "ɔ",
		49471: "ð",
		49472: "Á",
		49473: "Í",
		49474: "Ú",
		49475: "É",
		49476: "Ó",
		49477: "À",
		49478: "È",
		49479: "Ò",
		49480: "ì",
		49481: "ù",
		49482: "ý",
		49483: "ỳ",
		49484: "ɑ",
		49485: "ə",
		49487: "ɛ",
		49488: "θ",
		49489: "ʌ",
		49490: "ɑ́",
		49491: "ə́",
		49492: "ɔ́",
		49493: "ɛ́",
		49494: "ʌ́",
		49495: "ɑ̀",
		49496: "ə̀",
		49497: "ɔ̀",
		49498: "ɛ̀",
		49499: "ʌ̀",
		49500: "æ",
		49501: "ǽ",
		49502: "æ̀",
		49503: "Æ",
		49504: "ɑ̃",
		49505: "å",
		49506: "˘",
		49507: "ă",
		49508: "ŏ",
		49509: "ĭ",
		49510: "V́",
		49511: "T́",
		49513: "ɔ̃",
		49527: "ć",
		49531: "û",
		49532: "Ý",
		49534: "Ḿ",
		49700: "ō",
		49701: "ğ",
		49705: "Ḍ",
		49710: "Ḥ",
		49717: "Ṛ",
		49719: "Ṣ",
		49722: "Ẓ",
		49724: "ą",
		49728: "ḍ",
		49730: "ę",
		49734: "ḥ",
		49736: "ị",
		49740: "ṃ",
		49742: "ṇ",
		49747: "ṛ",
		49749: "ş",
		49750: "ṣ",
		49752: "ṭ",
		49757: "ẓ",
		49758: "İ",
		49759: "ṁ",
		49760: "ṅ",
		49761: "ż",
		49762: "Ś",
		49763: "ć",
		49764: "ń",
		49765: "ś",
		49766: "ý",
		49767: "ź",
		49768: "ì",
		49769: "Ä",
		49770: "Ö",
		49771: "Ü",
		49772: "ÿ",
		49773: "Â",
		49775: "û",
		49776: "Ā",
		49777: "Ē",
		49778: "Ī",
		49779: "Ō",
		49780: "Ū",
		49781: "ī",
		49782: "n̄",
		49783: "p̄",
		49784: "ū",
		49785: "ȳ",
		49786: "Ł",
		49787: "ł",
		49788: "ø",
		49789: "ĩ",
		49790: "õ",
		49955: "º",
		49956: "½",
		49958: "¹",
		49959: "²",
		49960: "¾",
		49961: "³",
		49972: "ɟ",
		50010: "g̀",
		50027: "ĕ",
		50028: "Č",
		50029: "Š",
		50030: "ǎ",
		50031: "č",
		50032: "ě",
		50033: "ň",
		50034: "ř",
		50035: "š",
		50036: "ž",
		50037: "ヰ",
		50038: "ヱ",
		50039: "ɯ̈",
		50040: "ɰ",
		50042: "ʔ",
		50043: "ɦ",
		50044: "ß",
		50209: "ɲ",
		50210: "ː",
	}
}

func (*daijirinExtractor) getFontWide() map[int]string {
	return map[int]string{
		41249: "仿",
		41250: "佉",
		41251: "侗",
		41252: "倘",
		41253: "偓",
		41254: "傔",
		41255: "傖",
		41256: "僄",
		41257: "僦",
		41258: "兕",
		41259: "凴",
		41260: "刁",
		41261: "剉",
		41262: "剗",
		41263: "劂",
		41264: "劓",
		41265: "勖",
		41266: "卬",
		41267: "厓",
		41268: "厲",
		41269: "呍",
		41270: "吧",
		41271: "咜",
		41272: "呫",
		41273: "呦",
		41274: "咿",
		41275: "咩",
		41276: "哿",
		41277: "唫",
		41278: "嘈",
		41279: "嘻",
		41280: "噯",
		41281: "噲",
		41282: "嚚",
		41283: "嚬",
		41284: "圊",
		41285: "圯",
		41286: "坌",
		41287: "埸",
		41288: "埶",
		41289: "埤",
		41290: "壔",
		41291: "壠",
		41292: "壚",
		41293: "虁",
		41294: "奝",
		41295: "奭",
		41296: "姒",
		41297: "婥",
		41298: "婕",
		41299: "孼",
		41300: "尫",
		41301: "屩",
		41302: "崧",
		41303: "嵆",
		41304: "嶠",
		41305: "嶸",
		41306: "幘",
		41307: "庾",
		41308: "龐",
		41309: "弇",
		41310: "彀",
		41311: "彐",
		41312: "彤",
		41313: "徉",
		41314: "徜",
		41315: "徧",
		41316: "忉",
		41317: "忼",
		41318: "忡",
		41319: "怵",
		41320: "悝",
		41321: "惛",
		41322: "惕",
		41323: "惙",
		41324: "惲",
		41325: "愷",
		41326: "戕",
		41327: "扃",
		41328: "扑",
		41329: "拖",
		41330: "拄",
		41331: "捃",
		41332: "挹",
		41333: "摹",
		41334: "撝",
		41335: "撿",
		41336: "昱",
		41337: "晡",
		41338: "皙",
		41339: "腊",
		41340: "臏",
		41341: "杇",
		41342: "枘",
		41505: "杻",
		41506: "棰",
		41507: "棖",
		41508: "楨",
		41509: "楣",
		41510: "橛",
		41511: "櫬",
		41512: "欛",
		41513: "歆",
		41514: "殂",
		41515: "殭",
		41516: "毱",
		41517: "氅",
		41518: "氐",
		41519: "氳",
		41520: "淼",
		41521: "沅",
		41522: "沆",
		41523: "汴",
		41524: "沔",
		41525: "泫",
		41526: "泮",
		41527: "洄",
		41528: "洎",
		41529: "洮",
		41530: "浥",
		41531: "淄",
		41532: "涿",
		41533: "淝",
		41534: "湜",
		41535: "渧",
		41536: "滃",
		41537: "漪",
		41538: "漚",
		41539: "漳",
		41540: "澌",
		41541: "瀆",
		41542: "灝",
		41543: "灤",
		41544: "灎",
		41545: "炫",
		41546: "炷",
		41547: "焮",
		41548: "焠",
		41549: "煜",
		41550: "煇",
		41551: "煆",
		41552: "煨",
		41553: "熅",
		41554: "熒",
		41555: "熇",
		41556: "熳",
		41557: "燋",
		41558: "燁",
		41559: "燾",
		41560: "凞",
		41561: "牓",
		41562: "牕",
		41563: "牖",
		41564: "犍",
		41565: "犛",
		41566: "猨",
		41567: "獐",
		41568: "獷",
		41569: "獼",
		41570: "玕",
		41571: "珉",
		41572: "琦",
		41573: "琚",
		41574: "琨",
		41575: "璆",
		41576: "璉",
		41577: "璟",
		41578: "璣",
		41579: "璘",
		41580: "璨",
		41581: "璿",
		41582: "瓚",
		41583: "畎",
		41584: "痀",
		41585: "痤",
		41586: "瘖",
		41587: "瘭",
		41588: "皞",
		41589: "盎",
		41590: "盌",
		41591: "盬",
		41592: "盼",
		41593: "眚",
		41594: "眙",
		41595: "睢",
		41596: "睟",
		41597: "睜",
		41598: "睽",
		41761: "矰",
		41762: "矻",
		41763: "砭",
		41764: "确",
		41765: "磈",
		41766: "磷",
		41767: "禘",
		41768: "秔",
		41769: "窅",
		41770: "窠",
		41771: "窬",
		41772: "窳",
		41773: "竽",
		41774: "筠",
		41775: "簋",
		41776: "簠",
		41777: "籮",
		41778: "糗",
		41779: "糕",
		41780: "糝",
		41781: "紈",
		41782: "紓",
		41783: "絇",
		41784: "絓",
		41785: "絜",
		41786: "絺",
		41787: "綈",
		41788: "緂",
		41789: "縈",
		41790: "縕",
		41791: "縑",
		41792: "縠",
		41793: "縝",
		41794: "繇",
		41795: "繒",
		41796: "繳",
		41797: "罽",
		41798: "罾",
		41799: "翟",
		41800: "翬",
		41801: "耦",
		41802: "聱",
		41803: "艴",
		41804: "芎",
		41805: "芷",
		41806: "芮",
		41807: "苾",
		41808: "茀",
		41809: "荇",
		41810: "荃",
		41811: "莘",
		41812: "蒯",
		41813: "蓰",
		41814: "蕓",
		41815: "蕙",
		41816: "蕞",
		41817: "蕤",
		41818: "薏",
		41819: "藿",
		41820: "蘐",
		41821: "虗",
		41822: "虢",
		41823: "虬",
		41824: "虯",
		41825: "虺",
		41826: "蚑",
		41827: "蚱",
		41828: "蜋",
		41829: "蝘",
		41830: "蝥",
		41831: "螈",
		41832: "螭",
		41833: "蠲",
		41834: "裊",
		41835: "裛",
		41836: "褰",
		41837: "袪",
		41838: "裎",
		41839: "裱",
		41840: "褚",
		41841: "觔",
		41842: "觖",
		41843: "觳",
		41844: "訕",
		41845: "訢",
		41846: "詘",
		41847: "詡",
		41848: "詹",
		41849: "誾",
		41850: "豨",
		41851: "豳",
		41852: "貒",
		41853: "賙",
		41854: "贛",
		42017: "跎",
		42018: "跗",
		42019: "踠",
		42020: "踔",
		42021: "踽",
		42022: "蹢",
		42023: "輞",
		42024: "輭",
		42025: "輶",
		42026: "轔",
		42027: "辧",
		42028: "辵",
		42029: "辶",
		42030: "辶",
		42031: "迤",
		42032: "邅",
		42033: "邈",
		42034: "邛",
		42035: "邢",
		42036: "邳",
		42037: "郅",
		42038: "鄧",
		42039: "鄱",
		42040: "鄴",
		42041: "酈",
		42042: "酛",
		42043: "酤",
		42044: "酴",
		42045: "醃",
		42046: "醞",
		42047: "醮",
		42048: "釃",
		42049: "釗",
		42050: "鈐",
		42051: "鈇",
		42052: "鉏",
		42053: "鉸",
		42054: "銈",
		42055: "鍈",
		42056: "鏜",
		42057: "鐲",
		42058: "鑊",
		42059: "鑣",
		42060: "閒",
		42061: "閟",
		42062: "閩",
		42063: "閽",
		42064: "闓",
		42065: "闐",
		42066: "闚",
		42067: "闞",
		42068: "阼",
		42069: "陘",
		42070: "隄",
		42071: "雒",
		42072: "雞",
		42073: "雩",
		42074: "靛",
		42075: "靳",
		42076: "鞺",
		42077: "韞",
		42078: "韛",
		42079: "韡",
		42080: "頫",
		42081: "顒",
		42082: "顓",
		42083: "顗",
		42084: "顥",
		42085: "颺",
		42086: "飥",
		42087: "餖",
		42088: "餼",
		42089: "餻",
		42090: "饘",
		42091: "駔",
		42092: "駙",
		42093: "騃",
		42094: "騶",
		42095: "騸",
		42096: "魞",
		42097: "鮏",
		42098: "鯁",
		42099: "鰶",
		42100: "鴞",
		42101: "鵷",
		42102: "鵰",
		42103: "鷃",
		42104: "麨",
		42105: "麼",
		42106: "黧",
		42107: "鼂",
		42108: "鼯",
		42109: "齁",
		42110: "齗",
		42273: "龔",
		42274: "捥",
		42275: "楤",
		42276: "丰",
		42278: "挊",
		42279: "艜",
		42280: "桒",
		42283: "亍",
		42284: "亹",
		42285: "儞",
		42286: "偁",
		42287: "儃",
		42288: "佪",
		42289: "儋",
		42290: "儈",
		42291: "侒",
		42292: "佷",
		42293: "伋",
		42294: "傜",
		42295: "淸",
		42296: "卺",
		42297: "划",
		42298: "勑",
		42299: "匇",
		42300: "匃",
		42301: "匜",
		42303: "嗢",
		42304: "囉",
		42305: "唽",
		42306: "嚕",
		42307: "噱",
		42308: "嘽",
		42309: "嚞",
		42310: "喁",
		42311: "噞",
		42313: "哯",
		42314: "嚩",
		42315: "喈",
		42317: "晷",
		42318: "叵",
		42319: "嗩",
		42320: "妋",
		42321: "娭",
		42322: "嫚",
		42323: "嬗",
		42325: "娓",
		42326: "姞",
		42328: "孁",
		42329: "堄",
		42330: "埿",
		42332: "坍",
		42333: "垸",
		42334: "坅",
		42335: "坷",
		42336: "壎",
		42337: "塤",
		42338: "堠",
		42339: "墪",
		42340: "埏",
		42341: "媳",
		42342: "墉",
		42343: "坨",
		42344: "圩",
		42345: "尰",
		42346: "屟",
		42347: "屣",
		42349: "异",
		42351: "岺",
		42352: "岏",
		42353: "巋",
		42354: "巑",
		42355: "帔",
		42356: "幉",
		42357: "帒",
		42358: "幞",
		42360: "彇",
		42361: "弣",
		42362: "弶",
		42363: "弽",
		42364: "庪",
		42365: "擌",
		42529: "擎",
		42530: "挗",
		42531: "擐",
		42532: "挍",
		42533: "搯",
		42534: "擷",
		42535: "掙",
		42536: "抳",
		42537: "攞",
		42538: "挃",
		42539: "撾",
		42540: "摭",
		42541: "熮",
		42543: "烑",
		42544: "灵",
		42545: "煑",
		42546: "爕",
		42547: "焄",
		42548: "獦",
		42549: "猧",
		42550: "猽",
		42551: "獒",
		42552: "獯",
		42553: "獫",
		42554: "玁",
		42555: "狁",
		42556: "狻",
		42557: "瀼",
		42558: "瀣",
		42559: "洿",
		42560: "濊",
		42561: "澠",
		42562: "潢",
		42563: "灊",
		42564: "淛",
		42565: "涘",
		42566: "湌",
		42567: "灔",
		42569: "涔",
		42570: "涬",
		42571: "邾",
		42572: "鄘",
		42573: "邶",
		42574: "鄀",
		42575: "鄽",
		42576: "菇",
		42577: "菆",
		42578: "蓀",
		42579: "藊",
		42580: "蘅",
		42581: "芺",
		42582: "蒺",
		42583: "蔾",
		42584: "蘼",
		42585: "薁",
		42586: "葒",
		42587: "蓯",
		42588: "蒾",
		42589: "蘩",
		42590: "蔌",
		42591: "蔞",
		42592: "菝",
		42593: "蕽",
		42594: "蘡",
		42595: "茛",
		42596: "荽",
		42597: "孽",
		42598: "葜",
		42599: "菀",
		42600: "薟",
		42601: "芾",
		42602: "蘘",
		42603: "蔲",
		42604: "蔯",
		42605: "荗",
		42606: "莔",
		42607: "噶",
		42608: "藋",
		42609: "莧",
		42610: "苆",
		42611: "蓪",
		42612: "萁",
		42613: "藦",
		42614: "薷",
		42615: "蘞",
		42616: "莕",
		42617: "蒅",
		42619: "芿",
		42620: "悆",
		42621: "忞",
		42622: "惸",
		42785: "惝",
		42786: "怳",
		42787: "惔",
		42788: "怍",
		42789: "惋",
		42790: "扆",
		42791: "曛",
		42792: "昀",
		42793: "昪",
		42794: "暍",
		42795: "臗",
		42796: "臛",
		42797: "膘",
		42798: "榺",
		42799: "樾",
		42800: "櫆",
		42801: "柀",
		42802: "棱",
		42803: "橒",
		42804: "檞",
		42805: "檨",
		42806: "杮",
		42807: "楉",
		42808: "樻",
		42810: "桕",
		42811: "棼",
		42812: "槾",
		42813: "楗",
		42814: "棙",
		42816: "桄",
		42817: "杴",
		42818: "枒",
		42819: "檫",
		42820: "杈",
		42821: "欋",
		42822: "棅",
		42823: "榀",
		42824: "棻",
		42825: "栭",
		42826: "榭",
		42827: "棌",
		42828: "欵",
		42829: "殩",
		42830: "殮",
		42831: "槩",
		42832: "櫲",
		42835: "穀",
		42836: "蒁",
		42837: "迱",
		42839: "适",
		42840: "逈",
		42841: "迍",
		42842: "逭",
		42843: "迮",
		42844: "璈",
		42845: "瑄",
		42846: "璱",
		42847: "玦",
		42848: "琯",
		42849: "璙",
		42850: "珅",
		42851: "珣",
		42852: "玠",
		42853: "瓈",
		42854: "璫",
		42855: "琫",
		42856: "瑍",
		42857: "琊",
		42858: "疿",
		42859: "癕",
		42860: "皥",
		42861: "皪",
		42862: "盦",
		42863: "盔",
		42864: "瞔",
		42865: "睠",
		42867: "瞟",
		42868: "瞍",
		42869: "眶",
		42871: "畾",
		42872: "矪",
		42873: "矬",
		42874: "穭",
		42876: "袽",
		42877: "襅",
		42878: "筯",
		43041: "帘",
		43042: "笇",
		43043: "篗",
		43044: "籡",
		43045: "籗",
		43046: "褲",
		43047: "褙",
		43048: "粿",
		43051: "縬",
		43052: "罇",
		43053: "纆",
		43054: "耖",
		43055: "耟",
		43056: "艉",
		43057: "賾",
		43058: "蟫",
		43059: "蜺",
		43060: "蚨",
		43061: "蟭",
		43062: "蠐",
		43063: "螬",
		43064: "蜟",
		43065: "蠼",
		43066: "螋",
		43067: "蚍",
		43068: "蟟",
		43069: "蛁",
		43070: "蜞",
		43073: "蝯",
		43075: "鵒",
		43076: "鴝",
		43077: "鸜",
		43078: "鸇",
		43079: "鶖",
		43081: "鸍",
		43082: "鵩",
		43083: "鶡",
		43084: "鷴",
		43086: "鷧",
		43087: "鏌",
		43088: "鎁",
		43089: "鍱",
		43090: "銙",
		43091: "釭",
		43092: "鉧",
		43093: "鍑",
		43094: "鏽",
		43095: "錕",
		43096: "鋂",
		43097: "鋧",
		43098: "鐴",
		43100: "鋐",
		43101: "蹔",
		43103: "踶",
		43104: "詵",
		43105: "諐",
		43106: "誮",
		43107: "謭",
		43108: "誷",
		43109: "觶",
		43110: "釄",
		43111: "醼",
		43112: "醨",
		43113: "釱",
		43114: "釻",
		43115: "鎛",
		43116: "鐧",
		43118: "鉃",
		43119: "纇",
		43120: "熲",
		43121: "頞",
		43122: "顖",
		43123: "蒴",
		43124: "蕺",
		43125: "芩",
		43126: "佺",
		43127: "佾",
		43128: "俏",
		43129: "倻",
		43130: "儵",
		43131: "噦",
		43132: "嗉",
		43133: "嘰",
		43134: "吒",
		43297: "唵",
		43298: "唼",
		43299: "埦",
		43300: "墝",
		43301: "埵",
		43302: "垜",
		43303: "墩",
		43304: "圳",
		43305: "壒",
		43306: "羗",
		43307: "搢",
		43308: "搩",
		43309: "攩",
		43310: "擤",
		43311: "挵",
		43312: "拼",
		43313: "擻",
		43314: "掽",
		43315: "湑",
		43316: "濹",
		43317: "泔",
		43318: "犎",
		43319: "桛",
		43320: "梣",
		43321: "樏",
		43322: "梻",
		43323: "橐",
		43324: "梘",
		43325: "梲",
		43326: "橅",
		43327: "檉",
		43329: "櫧",
		43330: "枻",
		43331: "柃",
		43332: "栱",
		43333: "栬",
		43334: "樝",
		43335: "橖",
		43336: "朳",
		43337: "棭",
		43338: "梂",
		43340: "榰",
		43341: "柷",
		43342: "槵",
		43343: "檔",
		43344: "桫",
		43345: "欏",
		43346: "枓",
		43347: "楲",
		43348: "腭",
		43349: "胳",
		43350: "腨",
		43351: "朓",
		43352: "鰧",
		43353: "蓏",
		43354: "玫",
		43355: "琰",
		43356: "瑇",
		43357: "璩",
		43358: "珧",
		43359: "瑀",
		43360: "瑒",
		43361: "瑭",
		43362: "玔",
		43363: "珖",
		43364: "玢",
		43365: "皶",
		43366: "麬",
		43367: "硨",
		43368: "磠",
		43369: "磤",
		43370: "磲",
		43371: "砍",
		43372: "硾",
		43373: "碰",
		43374: "硇",
		43375: "礀",
		43376: "畺",
		43377: "裰",
		43378: "裑",
		43379: "袘",
		43380: "襀",
		43381: "裓",
		43383: "褘",
		43384: "褹",
		43385: "襢",
		43386: "褨",
		43387: "篊",
		43388: "笧",
		43389: "簁",
		43390: "簎",
		43553: "簶",
		43554: "籰",
		43555: "籙",
		43556: "籭",
		43557: "箯",
		43558: "籑",
		43559: "荇",
		43560: "蓎",
		43561: "笯",
		43563: "篅",
		43564: "簳",
		43565: "簹",
		43566: "篔",
		43569: "筲",
		43570: "笭",
		43571: "筎",
		43572: "羖",
		43573: "籹",
		43574: "粏",
		43575: "糈",
		43576: "糫",
		43577: "粼",
		43578: "粔",
		43579: "粶",
		43580: "糙",
		43581: "糄",
		43582: "粬",
		43583: "糵",
		43584: "紽",
		43585: "緌",
		43586: "絁",
		43587: "紇",
		43588: "纑",
		43589: "緦",
		43590: "紞",
		43591: "纍",
		43593: "羿",
		43594: "翺",
		43595: "翥",
		43596: "羕",
		43597: "蝲",
		43598: "蟖",
		43599: "蚸",
		43600: "蜓",
		43601: "蜾",
		43602: "螇",
		43603: "蠁",
		43604: "蜱",
		43606: "蛺",
		43607: "虵",
		43608: "蝱",
		43609: "蠔",
		43610: "蝤",
		43611: "蛑",
		43612: "蠊",
		43613: "蠆",
		43614: "螠",
		43615: "鈸",
		43616: "錑",
		43617: "鎺",
		43618: "鍰",
		43619: "鏁",
		43620: "銲",
		43621: "鈹",
		43622: "鏟",
		43623: "鐖",
		43624: "鑯",
		43625: "闋",
		43627: "鏱",
		43628: "鈼",
		43630: "鬌",
		43631: "鞖",
		43632: "靪",
		43633: "鞚",
		43634: "靮",
		43635: "鬠",
		43636: "鱘",
		43637: "鮬",
		43638: "鱰",
		43639: "鱪",
		43640: "鯳",
		43641: "鱵",
		43642: "鯯",
		43643: "鯧",
		43644: "魳",
		43645: "鯎",
		43646: "鯥",
		43809: "鮄",
		43810: "鱩",
		43811: "鱮",
		43812: "鯇",
		43813: "鮞",
		43814: "鰖",
		43815: "鮸",
		43816: "鯷",
		43817: "魬",
		43818: "鯘",
		43819: "鱫",
		43820: "鱝",
		43821: "鱏",
		43822: "鱓",
		43823: "鰱",
		43824: "鮊",
		43825: "鱛",
		43826: "鮾",
		43827: "鱁",
		43828: "鮧",
		43829: "魦",
		43830: "鱭",
		43831: "孒",
		43832: "甪",
		43833: "厴",
		43834: "尩",
		43835: "车",
		43836: "电",
		43837: "邌",
		43838: "仐",
		43839: "么",
		43840: "蠃",
		43841: "兗",
		43842: "矠",
		43843: "矟",
		43844: "劻",
		43845: "勰",
		43846: "斲",
		43847: "姧",
		43848: "嬥",
		43849: "妤",
		43850: "媞",
		43852: "廋",
		43853: "庿",
		43854: "愒",
		43855: "憍",
		43856: "愐",
		43857: "豇",
		43858: "豉",
		43859: "雘",
		43860: "彔",
		43861: "邕",
		43862: "隺",
		43863: "幫",
		43864: "帮",
		43865: "毈",
		43867: "彽",
		43868: "徸",
		43869: "鄯",
		43870: "郄",
		43871: "邙",
		43872: "隩",
		43873: "犰",
		43874: "狳",
		43875: "獱",
		43876: "貛",
		43877: "攲",
		43878: "爗",
		43879: "滎",
		43880: "煠",
		43881: "燄",
		43882: "炻",
		43883: "烤",
		43884: "炗",
		43885: "剡",
		43886: "昉",
		43887: "昰",
		43888: "甗",
		43890: "瓫",
		43892: "敔",
		43893: "忩",
		43894: "毿",
		43895: "瘵",
		43896: "痎",
		43897: "癋",
		43898: "疒",
		43899: "癤",
		43900: "癭",
		43901: "瘙",
		43902: "痟",
		44065: "痏",
		44066: "眴",
		44067: "睺",
		44068: "毗",
		44069: "翮",
		44071: "稭",
		44072: "稹",
		44073: "祆",
		44074: "禖",
		44075: "皁",
		44076: "皝",
		44077: "翃",
		44078: "舢",
		44079: "艠",
		44082: "趯",
		44083: "醶",
		44084: "跑",
		44085: "蹰",
		44086: "躃",
		44087: "跆",
		44088: "韉",
		44089: "饠",
		44090: "躻",
		44091: "髹",
		44092: "髁",
		44093: "餛",
		44094: "餺",
		44095: "飣",
		44096: "飰",
		44097: "饆",
		44098: "靏",
		44099: "閦",
		44100: "闈",
		44101: "顬",
		44102: "頊",
		44103: "骶",
		44104: "髐",
		44106: "鶍",
		44107: "鴲",
		44108: "鸕",
		44109: "鵼",
		44110: "鷀",
		44112: "鼹",
		44113: "鼷",
		44114: "髖",
		44116: "鸊",
		44117: "鷉",
		44118: "鵟",
		44119: "鷟",
		44120: "鵂",
		44121: "鶹",
		44122: "鴗",
		44123: "鷚",
		44124: "鵇",
		44125: "鶊",
		44126: "鶼",
		44127: "觫",
		44128: "觘",
		44129: "觿",
		44130: "剕",
		44131: "颸",
		44132: "飇",
		44133: "飈",
		44134: "贉",
		44135: "賖",
		44136: "赬",
		44137: "鼗",
		44138: "鼐",
		44139: "鼺",
		44140: "齝",
		44141: "齭",
		44142: "齵",
		44143: "龗",
		44144: "蓂",
		44145: "藎",
		44146: "葼",
		44147: "茼",
		44148: "藭",
		44149: "薼",
		44150: "菪",
		44151: "莩",
		44152: "蓽",
		44153: "苕",
		44154: "芡",
		44155: "茺",
		44157: "蔤",
		44321: "葈",
		44322: "你",
		44323: "儛",
		44325: "塼",
		44326: "坼",
		44327: "塌",
		44328: "垿",
		44329: "姮",
		44330: "媧",
		44331: "嬙",
		44332: "渲",
		44333: "洦",
		44334: "滇",
		44335: "潙",
		44336: "澶",
		44337: "涮",
		44338: "涪",
		44339: "啐",
		44340: "嚈",
		44341: "噠",
		44342: "弴",
		44343: "哆",
		44344: "嚳",
		44345: "洱",
		44346: "灃",
		44347: "濞",
		44348: "湉",
		44349: "泆",
		44350: "洹",
		44351: "昫",
		44352: "暠",
		44353: "昕",
		44354: "昺",
		44355: "桲",
		44356: "橉",
		44357: "窼",
		44358: "穇",
		44359: "秫",
		44360: "秭",
		44361: "禛",
		44362: "祜",
		44363: "祹",
		44364: "蜇",
		44365: "蛼",
		44366: "蚜",
		44367: "蚉",
		44368: "蛽",
		44370: "螵",
		44371: "蚇",
		44372: "螓",
		44373: "蜐",
		44374: "瘀",
		44376: "瘼",
		44378: "痱",
		44379: "癯",
		44380: "癁",
		44381: "礴",
		44382: "礜",
		44383: "砉",
		44384: "耷",
		44385: "耼",
		44387: "軑",
		44388: "轘",
		44389: "輀",
		44390: "魹",
		44391: "韴",
		44392: "鞲",
		44395: "鮲",
		44397: "鰘",
		44399: "鰙",
		44400: "鯝",
		44401: "鰣",
		44402: "鯽",
		44404: "魶",
		44405: "鰚",
		44406: "鱲",
		44407: "鱜",
		44409: "鱊",
		44410: "鱐",
		44411: "鱟",
		44412: "魣",
		44413: "魫",
		44414: "驎",
		44577: "麯",
		44578: "驌",
		44579: "騮",
		44580: "驊",
		44581: "駃",
		44582: "騠",
		44583: "駰",
		44584: "騭",
		44585: "麅",
		44586: "麞",
		44589: "乚",
		44591: "氵",
		44592: "艹",
		44593: "艹",
		44594: "扌",
		44595: "阝",
		44596: "犭",
		44597: "阝",
		44598: "刂",
		44600: "忄",
		44602: "耂",
		44603: "爫",
		44605: "灬",
		44607: "氺",
		44609: "罒",
		44610: "礻",
		44611: "衤",
		44612: "飠",
		44632: "©",
		44633: "♮",
		44639: "㊙",
		44640: "☞",
		44649: "®",
		44651: "Æ",
		44652: "æ",
		44654: "ﬂ",
		44656: "œ",
		44657: "∘",
		44660: "℧",
		44663: "©",
		44665: "㊜",
		44666: "〖",
		44667: "〗",
		45108: "☰",
		45109: "☷",
		45110: "☱",
		45111: "☲",
		45112: "☴",
		45113: "☵",
		45114: "☶",
		45120: "℉",
		45121: "〽",
		45122: "卍",
		45123: "♨",
		45124: "♠",
		45125: "♥",
		45130: "♩",
		45133: "❶",
		45134: "❷",
		45135: "❸",
		45136: "❹",
		45137: "❺",
		45138: "❻",
		45139: "❼",
		45140: "❽",
		45141: "❾",
		45142: "❿",
		45143: "⓫",
		45144: "⓬",
		45145: "⓭",
		45146: "⓮",
		45147: "⓯",
		45148: "⓰",
		45149: "⓱",
		45150: "⓲",
		45151: "⓳",
		45158: "ノ",
		45163: "ヰ",
		45175: "㏋",
	}
}
