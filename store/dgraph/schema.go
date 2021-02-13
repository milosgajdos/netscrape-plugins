package dgraph

// Schema is DQL schema
const Schema = `
	type Entity {
		xid
		name
		namespace
		created_at
		links
	}

	type Resource {
		name
		group
		version
		kind
		namespaced
	}

	xid: string @index(exact) .
	name: string @index(exact) .
	namespace: string @index(exact) .
	created_at : datetime @index(hour) .
	links: [uid] @count @reverse .
	group: string @index(exact) .
	version: string @index(exact) .
	kind: string @index(exact) .
	namespaced: bool .
`
