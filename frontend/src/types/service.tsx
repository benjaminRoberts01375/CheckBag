class Service {
	internal_address: string;
	external_address: string;

	constructor(internal_address: string, external_address: string) {
		this.internal_address = internal_address;
		this.external_address = external_address;
	}
}

export default Service;
