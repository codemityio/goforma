package badge

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func coverage(c *cli.Context) error {
	coverageFilePath := c.String("cov-file-path")
	minimumCoverage := c.Int64("minimum")

	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%w: unable to determine working directory: %w", errPkg, err)
	}

	covPath, err := sanitisePath(baseDir, coverageFilePath)
	if err != nil {
		return fmt.Errorf("%w: invalid coverage file path `%s`", err, coverageFilePath)
	}

	covFile, err := os.Open(covPath) // #nosec G304
	if err != nil {
		return fmt.Errorf("%w: unable to read coverage file `%s`: %w", errPkg, covPath, err)
	}

	defer func() {
		_ = covFile.Close()
	}()

	var total string

	scanner := bufio.NewScanner(covFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "total") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				total = strings.TrimSuffix(fields[len(fields)-1], "%")
			}
		}
	}

	if total == "" {
		return fmt.Errorf("%w: unable to find coverage total in `%s`", errPkg, coverageFilePath)
	}

	parsedCoverage, err := strconv.ParseInt(strings.Split(total, ".")[0], 10, 32)
	if err != nil {
		return fmt.Errorf("%w: unable to parse coverage total `%s`: %w", errPkg, total, err)
	}

	docPath, err := sanitisePath(baseDir, documentPath)
	if err != nil {
		return fmt.Errorf("%w: invalid document path `%s`", err, documentPath)
	}

	cc, err := os.ReadFile(docPath) // #nosec G304
	if err != nil {
		return fmt.Errorf("%w: unable to read document file `%s`: %w", errPkg, docPath, err)
	}

	newBadge := fmt.Sprintf(
		"![%s](https://img.shields.io/badge/Coverage-%d%%25-%s.svg?longCache=true&style=flat)",
		id, parsedCoverage, badgeColour(parsedCoverage),
	)

	rrr := regexp.MustCompile(fmt.Sprintf(`!\[%s.*`, id))
	newReadme := rrr.ReplaceAll(cc, []byte(newBadge))

	if e := os.WriteFile(docPath, newReadme, perm); e != nil { // #nosec G703
		return fmt.Errorf("%w: unable to write document file `%s`: %w", errPkg, docPath, e)
	}

	if parsedCoverage < minimumCoverage {
		return fmt.Errorf("%w: coverage is %d < %d", errPkg, parsedCoverage, minimumCoverage)
	}

	return nil
}

func sanitisePath(baseDir, p string) (string, error) {
	if filepath.IsAbs(p) {
		return "", errPathOutsideBase
	}

	cleanPath := filepath.Clean(p)
	absPath := filepath.Join(baseDir, cleanPath)

	rel, err := filepath.Rel(baseDir, absPath)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errRelPath, err)
	}

	if strings.HasPrefix(rel, "..") {
		return "", errPathOutsideBase
	}

	return absPath, nil
}

func badgeColour(t int64) string {
	switch {
	case t < colourRed:
		return "red"
	case t < colourYellow:
		return "yellow"
	case t < colourYellowGreen:
		return "yellowgreen"
	case t < colourGreen:
		return "green"
	default:
		return "brightgreen"
	}
}
