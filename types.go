package weavychat

type (
	AppType string
)

const (
	AppTypeChat  AppType = "chat"
	AppTypeFiles AppType = "files"
	AppTypePosts AppType = "posts"
)

// String returns the string value of OTPType
func (o AppType) String() string {
	return string(o)
}
