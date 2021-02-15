package dgraph

// Schema is DQL schema
const Schema = `
	type Entity {
		xid
		name
		namespace
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
	links: [uid] @count @reverse .
	created_at : datetime @index(hour) .
	group: string @index(exact) .
	version: string @index(exact) .
	kind: string @index(exact) .
	namespaced: bool .
`
