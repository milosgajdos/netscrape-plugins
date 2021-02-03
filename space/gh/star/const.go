package star

const (
	// version is GitHub API version
	version = "v3"

	// ownerRes is owner resource name
	ownerRes = "owner"
	// ownerGroup is owner resource group
	ownerGroup = "owners"
	// ownerRol is the owner-repo relation
	ownerRel = "owns"

	// repoRes is repo resource name
	repoRes = "repo"
	// repoGroup is repo resource group
	repoGroup = "repos"

	// topicRes is topic resource name
	topicRes = "topic"
	// topicGroup is topic resource group
	topicGroup = "topics"
	// topicRel is the repo-topic relation
	topicRel = "hasTopic"

	// langRes is lang resource name
	langRes = "lang"
	// langGroup is lang resource group
	langGroup = "langs"
	// langRel is the repo-language relation
	langRel = "isLang"

	// resKind is the default space.Resource kind
	resKind = "starred"
	// ns is the name of the namespace
	ns = "global"
	// dateTime specifies date format
	dateTime = "2006-01-02T15:04:05"
)
