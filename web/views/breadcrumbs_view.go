package views

type BreadcrumbsLink struct {
	Href string
	Name string
}

type BreadcrumbsView struct {
	Items   []BreadcrumbsLink
	Options []BreadcrumbsLink
}
