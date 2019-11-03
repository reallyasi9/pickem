package pickem

// NameStringer is an interface that allows querying of a name of something.  Names come in two flavors:
// Name() - a way of representing the object in a pretty way that people can understand in a sentence.
// ShortName() - an abbreviation that takes up much less space than Name() for printing a lot of these things.
type NameStringer interface {
	Name() string
	ShortName() string
}
