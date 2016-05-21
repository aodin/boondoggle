package boondoggle

// Section is used to build the table of contents. It is a tree.
type Section struct {
	Name     string
	Children []Section
}
