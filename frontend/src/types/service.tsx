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
		if (internal_address.substring(0, 4) !== "http") {
			internal_address = "http://" + internal_address;
		}

		// Preserve hostname and port for internal address
		const internalUrl = new URL(internal_address);
		this.internal_address = internalUrl.port
			? `${internalUrl.hostname}:${internalUrl.port}`
			: internalUrl.hostname;

		// Handle external addresses similarly
		external_address.forEach((address, index) => {
			if (address.substring(0, 4) !== "http") {
				address = "http://" + address;
			}
			const externalUrl = new URL(address);
			external_address[index] = externalUrl.port
				? `${externalUrl.hostname}:${externalUrl.port}`
				: externalUrl.hostname;
		});

		this.external_address = external_address;
		this.title = title;
		this.id = id;
		this.clientID = crypto.randomUUID();
	}
}

export default Service;
