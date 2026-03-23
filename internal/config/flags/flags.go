package flags

const (
	Repository = iota
	Revision
	OrderBy
	UseCommitter
	Format
	Extensions
	Languages
	Exclude
	RestrictTo

	FlagCount

	BooleanFlag
	StringFlag
	StringSliceFlag
)

var (
	flagToName = map[int]string{
		Repository:   "repository",
		Revision:     "revision",
		OrderBy:      "order-by",
		UseCommitter: "use-committer",
		Format:       "format",
		Extensions:   "extensions",
		Languages:    "languages",
		Exclude:      "exclude",
		RestrictTo:   "restrict-to",
	}
	flagToUsage = map[int]string{
		Repository:   "путь до Git репозитория; по умолчанию текущая директория",
		Revision:     "указатель на коммит; HEAD по умолчанию",
		OrderBy:      "ключ сортировки результатов; один из lines (дефолт), commits, files",
		UseCommitter: "булев флаг, заменяющий в расчётах автора (дефолт) на коммиттера",
		Format:       "формат вывода; один из tabular (дефолт), csv, json, json-lines",
		Extensions:   "список расширений, сужающий список файлов в расчёте; множество ограничений разделяется запятыми",
		Languages:    "список языков (программирования, разметки и др.), сужающий список файлов в расчёте; множество ограничений разделяется запятыми",
		Exclude:      "набор Glob паттернов, исключающих файлы из расчёта",
		RestrictTo:   "набор Glob паттернов, исключающий все файлы, не удовлетворяющие ни одному из паттернов набора",
	}
	flagToValueType = map[int]int{
		Repository:   StringFlag,
		Revision:     StringFlag,
		OrderBy:      StringFlag,
		UseCommitter: BooleanFlag,
		Format:       StringFlag,
		Extensions:   StringSliceFlag,
		Languages:    StringSliceFlag,
		Exclude:      StringSliceFlag,
		RestrictTo:   StringSliceFlag,
	}
	flagToDefaultValue = map[int]interface{}{
		Repository:   ".",
		Revision:     "HEAD",
		OrderBy:      "lines",
		UseCommitter: false,
		Format:       "tabular",
		Extensions:   []string{},
		Languages:    []string{},
		Exclude:      []string{},
		RestrictTo:   []string{},
	}
)

func FlagDefaultValue(flag int) any {
	return flagToDefaultValue[flag]
}

func createValue(valueType, flagType int) Value {
	switch valueType {
	case BooleanFlag:
		return newBoolValue(flagToDefaultValue[flagType].(bool))
	case StringFlag:
		return newStringValue(flagToDefaultValue[flagType].(string))
	case StringSliceFlag:
		return newStringSliceValue(flagToDefaultValue[flagType].([]string))
	default:
		panic("unhandled default case")
	}
}

func CreateFlags() map[int]*Flag {
	flags := make(map[int]*Flag, FlagCount)
	for flagIota := 0; flagIota < FlagCount; flagIota++ {
		flags[flagIota] = &Flag{
			Name:  flagToName[flagIota],
			Use:   flagToUsage[flagIota],
			Value: createValue(flagToValueType[flagIota], flagIota),
		}
	}
	return flags
}
