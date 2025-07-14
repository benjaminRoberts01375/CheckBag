class Service {
	internal_address: string;
	external_address: string[];
	title: string;
	id: string;
	clientID: string;

	constructor(
		internal_address: string,
		external_address: string[],
		title: string,
		id: string = "",
	) {
		this.internal_address = internal_address;
		external_address.forEach((address, index) => {
			external_address[index] = new URL(address).hostname;
		});

		this.external_address = external_address;
		this.title = title;
		this.id = id;
		this.clientID = crypto.randomUUID();
	}
}

export default Service;
