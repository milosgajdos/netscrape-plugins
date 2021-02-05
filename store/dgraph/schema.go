package dgraph

// Schema is dgraph DQL schema
const Schema = `
	type Entity {
		xid
		name
		kind
		namespace
		created_at
		link
	}

	xid: string @index(exact) .
	name: string @index(exact) .
	kind: string @index(exact) .
	namespace: string @index(exact) .
	created_at : datetime @index(hour) .
	link: [uid] @count @reverse .
`
