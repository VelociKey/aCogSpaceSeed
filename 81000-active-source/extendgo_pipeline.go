package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"sov.fleet/extendgo/81000-active-source/emitter"
	"sov.fleet/extendgo/81000-active-source/ghost"
	"sov.fleet/extendgo/81000-active-source/overlay"
	"sov.fleet/extendgo/81000-active-source/spec_ingest"
	"sov.fleet/extendgo/81000-active-source/verifier"
)

func main() {
	specFlag := flag.String("spec", "", "Path to Go language specification HTML/text file")
	outFlag := flag.String("out", "", "Target output path for generated grammar_go.wag")
	deltaOutFlag := flag.String("delta-out", "", "Target output path for generated version delta extend_go_v<version>.webnf")
	versionFlag := flag.String("version", "1.27", "Target Go version")
	verifyOnlyFlag := flag.Bool("verify", false, "Verify without writing to destination")
	auditFlag := flag.String("audit-visibility", "", "Path to Go source file to audit for ghost vs physical visibility")
	modernizeFlag := flag.String("modernize", "", "Path to Go source file to modernize with materialized explicit suffixes")
	modernizeOutFlag := flag.String("modernize-out", "", "Output path for modernized Go source (optional)")

	flag.Parse()

	// Handle ghost visibility audit if requested
	if *auditFlag != "" {
		handleAuditVisibility(*auditFlag)
		return
	}

	// Handle ghost modernization if requested
	if *modernizeFlag != "" {
		handleModernizeSource(*modernizeFlag, *modernizeOutFlag)
		return
	}

	// 1. Resolve Spec Path
	specPath := *specFlag
	if specPath == "" {
		specPath = autoDiscoverGoSpec()
	}

	if specPath == "" {
		slog.Error("Failed to locate Go language specification file (doc/go_spec.html)")
		os.Exit(1)
	}

	// Detect version from spec content if possible
	detectedVersion := *versionFlag
	if specContent, err := os.ReadFile(specPath); err == nil {
		detectedVersion = overlay.DetectGoVersionFromSpec(string(specContent), *versionFlag)
	}

	slog.Info("Ingesting upstream Go specification", "path", specPath, "version", detectedVersion)

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

	// 4. Apply Sovereign Extension Overlay Functor & Perform Subsumption Audit
	overlayConfig := overlay.DefaultGoExtensions()
	overlayConfig.TargetVersion = detectedVersion
	functor := overlay.NewRuleOverlayFunctor(overlayConfig)
	synthesizedRules, err := functor.Apply(rawRules)
	if err != nil {
		slog.Error("Rule overlay synthesis failed", "error", err)
		os.Exit(1)
	}

	if rep := functor.LastSubsumptionReport; rep != nil {
		slog.Info("Subsumption & convergence audit complete",
			"go_version", rep.GoVersion,
			"total_extensions", rep.TotalExtensions,
			"active_extensions", rep.ActiveCount,
			"subsumed_upstream", rep.SubsumedCount,
		)
		for _, sf := range rep.SubsumedFeatures {
			slog.Info("[SUBSUMED] Feature natively converged into Go",
				"feature", sf.Feature,
				"subsumed_in", sf.SubsumedIn,
				"rule", sf.TargetRule,
				"reason", sf.Reason,
			)
		}
		for _, act := range rep.ActiveExtensions {
			slog.Info("[ACTIVE] Sovereign extension required",
				"name", act.Name,
				"action", act.Action,
				"target_rule", act.TargetRule,
			)
		}
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
	wagContent := emitter.EmitWebNFSN(synthesizedRules, detectedVersion)

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

	// 8. Emit version-specific realized extension program (.webnf)
	if functor.Spec != nil && functor.LastSubsumptionReport != nil {
		deltaWebNF := overlay.EmitVersionDeltaWebNF(functor.Spec, functor.LastSubsumptionReport)
		deltaPath := *deltaOutFlag
		if deltaPath == "" {
			deltaPath = filepath.Join(`C:\aCogSpaceSeed\00flow\extendgo\81000-active-source\overlay`, fmt.Sprintf("extend_go_v%s.webnf", detectedVersion))
		}
		if err := os.WriteFile(deltaPath, []byte(deltaWebNF), 0644); err != nil {
			slog.Warn("Failed to write version delta webnf", "path", deltaPath, "error", err)
		} else {
			slog.Info("Emitted realized version extension specification", "path", deltaPath, "bytes", len(deltaWebNF))
			invokeLinterGates(deltaPath)
		}
	}

	// 9. Invoke Platform Linter Gates
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

func handleAuditVisibility(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read file for visibility audit", "path", path, "error", err)
		os.Exit(1)
	}

	report, err := ghost.ScanAndMaterializeSource(data)
	if err != nil {
		slog.Error("Visibility audit failed", "path", path, "error", err)
		os.Exit(1)
	}

	slog.Info("=== SOVEREIGN GHOST VISIBILITY AUDIT ===",
		"file", path,
		"total_decls", report.TotalDecls,
		"ghost_public", report.GhostPublic,
		"ghost_private", report.GhostPrivate,
		"physical_public", report.PhysicalPublic,
		"physical_private", report.PhysicalPrivate,
	)

	for _, d := range report.Declarations {
		status := "GHOST"
		if !d.Visibility.IsGhost {
			status = "PHYSICAL"
		}
		fmt.Printf("  [%s] %-5s %-25s -> %s (casing: %s)\n",
			status, d.Kind, d.Name, d.Visibility.Kind, d.OriginalCasing)
	}
}

func handleModernizeSource(inPath, outPath string) {
	data, err := os.ReadFile(inPath)
	if err != nil {
		slog.Error("Failed to read file for modernization", "path", inPath, "error", err)
		os.Exit(1)
	}

	modernized, report, err := ghost.ModernizeSource(data)
	if err != nil {
		slog.Error("Modernization failed", "path", inPath, "error", err)
		os.Exit(1)
	}

	slog.Info("=== SOVEREIGN SOURCE MODERNIZATION COMPLETE ===",
		"file", inPath,
		"materialized_public", report.GhostPublic,
		"materialized_private", report.GhostPrivate,
	)

	if outPath == "" {
		fmt.Println(string(modernized))
	} else {
		if err := os.WriteFile(outPath, modernized, 0644); err != nil {
			slog.Error("Failed to write modernized source", "dest", outPath, "error", err)
			os.Exit(1)
		}
		slog.Info("Successfully wrote modernized source", "dest", outPath, "bytes", len(modernized))
	}
}
