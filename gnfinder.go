package main

import (
  "fmt"
  "flag"
  "encoding/csv"
  "os"
  "io"
  "strings"
  "github.com/gnames/gnfinder"
  "github.com/gnames/gnfinder/ent/nlp"
  "github.com/gnames/gnfinder/io/dict"
  "github.com/gnames/gnfinder/config"
  "github.com/gnames/gnfinder/ent/lang"
)

func main() {
  var colIndex int
  var bayesOddsThreshold float64

  flag.IntVar(&colIndex, "c", 1, "name index column")
  flag.Float64Var(&bayesOddsThreshold, "t", 80, "bayes odds threshold")
  flag.Parse()

  opts := []config.Option{
          config.OptLanguage(lang.English),
          config.OptBayesOddsThreshold(bayesOddsThreshold),
          config.OptWithAllMatches(true),
          config.OptWithAmbiguousNames(true),
          config.OptWithOddsAdjustment(true),

  }
  cfg := config.New(opts...)
  dictionary := dict.LoadDictionary()
  weights := nlp.BayesWeights()
  gnf := gnfinder.New(cfg, dictionary, weights)

  rawData, err := io.ReadAll(os.Stdin)
  if err != nil {
    fmt.Println("Can't read from stdin")
    os.Exit(1)
  }

  data := strings.ReplaceAll(string(rawData), "\"", "")
  cr := csv.NewReader(strings.NewReader(data))
  cw := csv.NewWriter(os.Stdout)
  cr.Comma = '\t'
  cw.Comma = '\t'
  cr.LazyQuotes = true
  for {
    rec, err := cr.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    txt := rec[colIndex]
    res := gnf.Find("", txt)
	  if len(res.Names) > 0 {
      for _, name := range res.Names {
        rec = append(rec, name.Name)
        cw.Write(rec)
      }
    } else {
      rec = append(rec, "NOT_FOUND")
      cw.Write(rec)
	  }
  }

  cw.Flush()
}

