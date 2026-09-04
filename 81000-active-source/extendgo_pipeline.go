package main

import (
	"flag"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"sov.fleet/extendgo/81000-active-source/emitter"
	"sov.fleet/extendgo/81000-active-source/overlay"
	"sov.fleet/extendgo/81000-active-source/spec_ingest"
	"sov.fleet/extendgo/81000-active-source/verifier"
)

func main() {
	specFlag := flag.String("spec", "", "Path to Go language specification HTML/text file")
	outFlag := flag.String("out", "", "Target output path for generated grammar_go.wag")
	versionFlag := flag.String("version", "1.27", "Target Go version")
	verifyOnlyFlag := flag.Bool("verify", false, "Verify without writing to destination")

	flag.Parse()

	// 1. Resolve Spec Path
	specPath := *specFlag
	if specPath == "" {
		specPath = autoDiscoverGoSpec()
	}

	if specPath == "" {
		slog.Error("Failed to locate Go language specification file (doc/go_spec.html)")
		os.Exit(1)
	}

	slog.Info("Ingesting upstream Go specification", "path", specPath, "version", *versionFlag)

	// 2. Ingest & Extract EBNF
	extracted, err := spec_ingest.ExtractEBNFFromFile(specPath)
	if err != nil {
		slog.Error("EBNF extraction failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Extracted EBNF rule blocks", "blocks", len(extracted.RuleBlocks))

	// 3. Parse EBNF into RawRules
	rawRules, err := spec_ingest.ParseEBNFRules(extracted.FullText)
	if err != nil {
		slog.Error("EBNF parsing failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Parsed base Go production rules", "count", len(rawRules))

	// 4. Apply Sovereign Extension Overlay Functor
	functor := overlay.NewRuleOverlayFunctor(overlay.DefaultGoExtensions())
	synthesizedRules, err := functor.Apply(rawRules)
	if err != nil {
		slog.Error("Rule overlay synthesis failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Applied sovereign rule overlay", "synthesized_rules", len(synthesizedRules))

	// 5. Verify Constraints & Sovereign Precision Audit
	report := verifier.VerifyGrammarConstraints(synthesizedRules)
	if !report.Valid {
		for _, v := range report.Violations {
			slog.Error("Grammar constraint violation", "issue", v)
		}
		os.Exit(1)
	}

	slog.Info("All grammar constraints & negative guards passed", "checked_rules", report.CheckedRulesCount)

	// 6. Emit strictly formatted webnf.sn grammar
	wagContent := emitter.EmitWebNFSN(synthesizedRules, *versionFlag)

	if *verifyOnlyFlag {
		slog.Info("Verification completed successfully (dry-run mode).")
		os.Exit(0)
	}

	// 7. Write to Destination
	outPath := *outFlag
	if outPath == "" {
		outPath = autoDiscoverOutputWAG()
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		slog.Error("Failed to create destination directory", "error", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outPath, []byte(wagContent), 0644); err != nil {
		slog.Error("Failed to write output grammar", "path", outPath, "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully synthesized and emitted grammar", "destination", outPath, "bytes", len(wagContent))

	// 8. Invoke Platform Linter Gates
	invokeLinterGates(outPath)
}

func autoDiscoverGoSpec() string {
	candidates := []string{
		`C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\go\doc\go_spec.html`,
		`C:\aCogSpaceSeed\86sref\golang-go\doc\go_spec.html`,
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

func autoDiscoverOutputWAG() string {
	return `C:\aCogSpaceSeed\00flow\languagegrammar\81000-active-source\grammar_go.wag`
}

func invokeLinterGates(targetWag string) {
	linters := []string{
		`C:\aCogSpaceSeed\00flow\latentlingua\webnf-sn-lint.exe`,
		`C:\aCogSpaceSeed\00flow\latentlingua\sn-lint.exe`,
	}

	for _, l := range linters {
		if _, err := os.Stat(l); err == nil {
			cmd := exec.Command(l, targetWag)
			out, err := cmd.CombinedOutput()
			if err != nil {
				slog.Warn("Linter gate check warning", "linter", filepath.Base(l), "error", err, "output", string(out))
			} else {
				slog.Info("Linter gate passed", "linter", filepath.Base(l))
			}
		}
	}
}
