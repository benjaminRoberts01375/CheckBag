class APIKey {
	name: string;
	id: string;
	key?: string;

	constructor(name: string, key_id: string, key?: string) {
		this.name = name;
		this.key = key;
		this.id = key_id;
	}
}

export default APIKey;
