package config

import (
	"fmt"
	"os"

	"gitlab.com/slon/shad-go/gitfame/internal/config/flags"
	"gitlab.com/slon/shad-go/gitfame/internal/errors"

	"github.com/spf13/pflag"
)

var (
	validOrderKeys = map[string]struct{}{
		"lines":   {},
		"commits": {},
		"files":   {},
	}
	validFormatKeys = map[string]struct{}{
		"tabular":    {},
		"csv":        {},
		"json":       {},
		"json-lines": {},
	}
)

type LanguageResolver interface {
	Resolve(language string) ([]string, bool)
}

type Config struct {
	Repository   string
	Revision     string
	OrderBy      string
	UseCommitter bool
	Format       string
	Extensions   []string
	Exclude      []string
	RestrictTo   []string
}

func MustLoad(resolver LanguageResolver) *Config {
	flagsMap := flags.CreateFlags()

	for i := 0; i < flags.FlagCount; i++ {
		f := flagsMap[i]
		pflag.VarP(f.Value, f.Name, "", f.Use)
		if f.Value.Type() == "bool" {
			pflag.Lookup(f.Name).NoOptDefVal = "true"
		}
	}

	pflag.Parse()

	rep, _ := flagsMap[flags.Repository].GetString()
	rev, _ := flagsMap[flags.Revision].GetString()
	order, _ := flagsMap[flags.OrderBy].GetString()
	uc, _ := flagsMap[flags.UseCommitter].GetBool()
	fmtStr, _ := flagsMap[flags.Format].GetString()
	exts, _ := flagsMap[flags.Extensions].GetStringSlice()
	langs, _ := flagsMap[flags.Languages].GetStringSlice()
	excl, _ := flagsMap[flags.Exclude].GetStringSlice()
	restr, _ := flagsMap[flags.RestrictTo].GetStringSlice()

	if _, allowed := validOrderKeys[order]; !allowed {
		fmt.Printf(errors.MsgInvalidOrderBy, order)
		os.Exit(1)
	}

	if _, allowed := validFormatKeys[fmtStr]; !allowed {
		fmt.Printf(errors.MsgInvalidFormat, fmtStr)
		os.Exit(1)
	}

	for _, lang := range langs {
		langExts, found := resolver.Resolve(lang)
		if !found {
			// fmt.Printf(errors.MsgUnknownLanguage, lang) // log error but not in stdout
			continue
		}

		exts = append(exts, langExts...)
	}

	return &Config{
		Repository:   rep,
		Revision:     rev,
		OrderBy:      order,
		UseCommitter: uc,
		Format:       fmtStr,
		Extensions:   exts,
		Exclude:      excl,
		RestrictTo:   restr,
	}
}
