package golib

import "regexp"

// --- pre-compiled patterns (compile once, reuse many times) ----
var (
	// characters that Windows and most POSIX file systems forbid
	patIllegal = regexp.MustCompile(`[/?<>\\:*|"]`)
	// special Windows device names such as CON, PRN, AUX, NUL, COM1 … LPT9
	patWinReserved = regexp.MustCompile(`(?i)^(con|prn|aux|nul|com[0-9]|lpt[0-9])(\..*)?$`)
	// ASCII control chars (0x00-0x1F) plus DEL-ish C1 controls (0x80-0x9F)
	patControlChars = regexp.MustCompile(`[\x00-\x1f\x80-\x9f]`)
	// a name that is just dots (“.”, “..”, “...”, …)
	patReserved = regexp.MustCompile(`^\.+$`)
	// leading dots or spaces
	patLeading = regexp.MustCompile(`^[\. ]+`)
	// trailing dots or spaces
	patTrailing = regexp.MustCompile(`[\. ]+$`)
)

// sanitizeBaseName scrubs path-unsafe characters and quirks, then
// (optionally) trims the result to ≤ 255 runes so it is safe on
// all common desktop and mobile file systems.
func sanitizeBaseName(
	fileName string,
	limitLength bool,
	removeUnsafe bool, // true  ➜ drop bad chars entirely
	replacementRune rune, // false ➜ replace with this rune (e.g. '_')
) string {

	replacement := string(replacementRune)
	if removeUnsafe {
		replacement = ""
	}

	s := fileName
	s = patIllegal.ReplaceAllString(s, replacement)
	s = patControlChars.ReplaceAllString(s, replacement)
	s = patReserved.ReplaceAllString(s, replacement)
	s = patWinReserved.ReplaceAllString(s, replacement)
	s = patLeading.ReplaceAllString(s, replacement)
	s = patTrailing.ReplaceAllString(s, replacement)

	if limitLength {
		r := []rune(s)
		if len(r) > 255 {
			s = string(r[:255]) // cut on rune boundary (avoids UTF-8 corruption)
		}
	}
	return s
}

// SanitizeBaseName mirrors the Java one-liner:
// • trims to 255 runes
// • keeps bad chars by swapping them for '_' (underscore)
func SanitizeBaseName(fileName string) string {
	return sanitizeBaseName(fileName, true, false, '_')
}
