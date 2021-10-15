package i18n

const (
	ErrUnexpected Message = iota + 1
	ErrBadRequest
	ErrSaveFailed
	ErrDeleteFailed
	ErrAlreadyExists
	ErrNotFound
	ErrSelectionNotFound
	ErrEntityNotFound
	ErrUserNotFound
	ErrUnauthorized
	ErrNoItemsSelected
	ErrConnectionFailed
	ErrInvalidCredentials
	ErrMissingParameter
	ErrProfanity

	MsgChangesSaved
	MsgEntryAddedTo
	MsgEntriesAddedTo
	MsgEntryRemovedFrom
	MsgEntriesRemovedFrom
	MsgEntryRemoved
	MsgUserDeleted
	MsgUserCreated
)

var Messages = MessageMap{
	// Error messages:
	ErrUnexpected:         gettext("Unexpected error, please try again"),
	ErrBadRequest:         gettext("Invalid request"),
	ErrSaveFailed:         gettext("Changes could not be saved"),
	ErrDeleteFailed:       gettext("Could not be deleted"),
	ErrAlreadyExists:      gettext("%s already exists"),
	ErrNotFound:           gettext("Not found on server, deleted?"),
	ErrSelectionNotFound:  gettext("Selection not found"),
	ErrEntityNotFound:     gettext("Not found on server, deleted?"),
	ErrUserNotFound:       gettext("User not found"),
	ErrUnauthorized:       gettext("Please log in and try again"),
	ErrNoItemsSelected:    gettext("No items selected"),
	ErrConnectionFailed:   gettext("Could not connect, please try again"),
	ErrInvalidCredentials: gettext("Invalid credentials"),
	ErrMissingParameter:   gettext("missing parameter %s"),
	ErrProfanity:          gettext("input contains profane text"),

	// Info and confirmation messages:
	MsgChangesSaved:       gettext("Changes successfully saved"),
	MsgEntryAddedTo:       gettext("One entry added to %s"),
	MsgEntriesAddedTo:     gettext("%d entries added to %s"),
	MsgEntryRemovedFrom:   gettext("One entry removed from %s"),
	MsgEntriesRemovedFrom: gettext("%d entries removed from %s"),
	MsgUserDeleted:        gettext("Account deleted"),
	MsgUserCreated:        gettext("User created"),
	MsgEntryRemoved:       gettext("Entry deleted successfully"),
}
