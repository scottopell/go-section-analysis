// Based heavily on https://stackoverflow.com/a/70777803
// Basically taken as-is and added some formatting options
package main

import (
	"debug/elf"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/go-delve/delve/pkg/proc"
)

type SectionMap map[string]int
type PkgSectionMap map[string]SectionMap

func doSort(pkgs []string, usage PkgSectionMap, sortBySection string) []string {
	sortedPkgs := make([]string, len(pkgs))

	sectionSizePerPkg := make(map[string]int)

	for _, pkg := range pkgs {
		sections, exists := usage[pkg]
		if !exists {
			continue
		}

		sectionFound := false
		for section, size := range sections {
			if section == sortBySection {
				sectionSizePerPkg[pkg] = size
				sectionFound = true
			}
		}
		if sectionFound {
			sortedPkgs = append(sortedPkgs, pkg)
		}
	}

	sort.Slice(sortedPkgs, func(i, j int) bool {
		iSectionSize := sectionSizePerPkg[sortedPkgs[i]]
		jSectionSize := sectionSizePerPkg[sortedPkgs[j]]
		return iSectionSize < jSectionSize
	})

	return sortedPkgs
}

func printFinalOutput(pkgs []string, usage PkgSectionMap) {
	for _, pkg := range pkgs {
		sections, exists := usage[pkg]
		if !exists {
			continue
		}

		fmt.Printf("%s:\n", pkg)
		for section, size := range sections {
			fmt.Printf("%20s: %8d bytes (%8s)\n", section, size, humanize.Bytes(uint64(size)))
		}
		fmt.Println()
	}
}

func printOverallUsage(usage SectionMap) {
	fmt.Printf("Total Section Usage:\n")
	for sectionName, size := range usage {

		fmt.Printf("%15s: %8d bytes (%8s)\n", sectionName, size, humanize.Bytes(uint64(size)))
	}
}

func main() {
	// Use delve to decode the DWARF section
	binInfo := proc.NewBinaryInfo(runtime.GOOS, runtime.GOARCH)
	err := binInfo.AddImage(os.Args[1], 0)
	if err != nil {
		panic(err)
	}

	// Make a list of unique packages
	pkgs := make([]string, 0, len(binInfo.PackageMap))
	for _, fullPkgs := range binInfo.PackageMap {
		for _, fullPkg := range fullPkgs {
			exists := false
			for _, pkg := range pkgs {
				if fullPkg == pkg {
					exists = true
					break
				}
			}
			if !exists {
				pkgs = append(pkgs, fullPkg)
			}
		}
	}
	// Sort them for a nice output
	sort.Strings(pkgs)

	// Parse the ELF file ourselfs
	elfFile, err := elf.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Get the symbol table
	symbols, err := elfFile.Symbols()
	if err != nil {
		panic(err)
	}

	perPkgUsage := make(PkgSectionMap)
	overallUsage := make(SectionMap)

	for _, sym := range symbols {
		if sym.Section == elf.SHN_UNDEF || sym.Section >= elf.SectionIndex(len(elfFile.Sections)) {
			continue
		}

		sectionName := elfFile.Sections[sym.Section].Name

		symPkg := ""
		for _, pkg := range pkgs {
			if strings.HasPrefix(sym.Name, pkg) {
				symPkg = pkg
				break
			}
		}

		overallUsage[sectionName] += int(sym.Size)

		// Symbol doesn't belong to a known package
		if symPkg == "" {
			continue
		}

		pkgStats := perPkgUsage[symPkg]
		if pkgStats == nil {
			pkgStats = make(map[string]int)
		}

		pkgStats[sectionName] += int(sym.Size)
		perPkgUsage[symPkg] = pkgStats
	}

	sortedPkgList := doSort(pkgs, perPkgUsage, ".rodata")

	printFinalOutput(sortedPkgList, perPkgUsage)
	printOverallUsage(overallUsage)
}
