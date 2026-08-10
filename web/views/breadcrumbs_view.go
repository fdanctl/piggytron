package views

// BreadcrumbsLink is one breadcrumb entry: its URL and label.
type BreadcrumbsLink struct {
	Href string
	Name string
}

// BreadcrumbsView holds the breadcrumb trail and, optionally, the sibling
// links shown as a dropdown at the current level.
type BreadcrumbsView struct {
	Items   []BreadcrumbsLink
	Options []BreadcrumbsLink
}
