package actions

type PageType string

func (pt PageType) String() string {
	return string(pt)
}

const (
	PAGE_ADMIN_LOGIN = PageType("/admin/login")
	PAGE_ADMIN_MAIN  = PageType("/admin")
	PAGE_NOTHING     = PageType("")
)
