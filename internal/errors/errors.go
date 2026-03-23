package errors

const (
	MsgLoadLanguageBase = "failed to load language base: %v\n"
	MsgGetTargetFiles   = "failed to get target files: %v\n"
	MsgWorkerBlame      = "failed on file %s: %w"
	MsgAggregateBlame   = "failed to blame file: %v\n"
	MsgUnknownOrderBy   = "unknown order-by parameter: %s\n"
	MsgPrintStats       = "failed to print stats: %v\n"

	MsgGitBlame  = "git blame failed on %s: %w"
	MsgGitLog    = "git log failed on %s: %w"
	MsgGitLsTree = "git ls-tree failed: %w"

	MsgUnknownFormat = "unknown format: %s"

	MsgInvalidOrderBy  = "invalid order-by: %s\n"
	MsgInvalidFormat   = "invalid format: %s\n"
	MsgUnknownLanguage = "warning: unknown language %q\n"

	MsgFlagTypeMismatch = "%w of %s flag. The value should be %s"
)
