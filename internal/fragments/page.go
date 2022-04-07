package fragments

type MenuItem struct {
	Icon    string
	Link    string
	Caption string
	Active  bool
}

type PageModel struct {
	Navigation    []MenuItem
	AppBar        []MenuItem
	AppBarCaption string
	Hamburger     []MenuItem
	Fragments     []FragmentModel
}
