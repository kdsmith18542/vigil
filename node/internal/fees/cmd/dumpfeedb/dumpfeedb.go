// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Tool dumpfeedb can be used to dump the internal state of the buckets of an
// estimator's feedb so that it can be externally analyzed.
package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/vigilnetwork/vgl/VGLutil/v4"
	"github.com/vigilnetwork/vgl/internal/fees"
	flags "github.com/jessevdk/go-flags"
)

type config struct {
	DB string `short:"b" long:"db" description:"Path to fee database"`
}

func main() {
	cfg := config{
		DB: path.Join(VGLutil.AppDataDir("vgld", false), "data", "mainnet", "feesdb"),
	}

	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		var e *flags.Error
		if !errors.As(err, &e) || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return
	}

	ecfg := fees.EstimatorConfig{
		DatabaseFile:         cfg.DB,
		ReplaceBucketsOnLoad: true,
		MinBucketFee:         1,
		MaxBucketFee:         2,
		FeeRateStep:          fees.DefaultFeeRateStep,
	}
	est, err := fees.NewEstimator(&ecfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(est.DumpBuckets())
}
