package dgraph

// SpaceSchema is space DQL schema
const SpaceSchema = `
	type Object {
		xid
		name
		namespace
		resource
		links
	}

	type Resource {
		xid
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
	resource: uid @count @reverse .
`
