class Analytic {
	quantity: number;
	country: Map<string, number>;
	ip: Map<string, number>;
	resource: Map<string, number>;

	constructor(
		quantity: number,
		country: Map<string, number>,
		ip: Map<string, number>,
		resource: Map<string, number>,
	) {
		this.quantity = quantity;
		this.country = country;
		this.ip = ip;
		this.resource = resource;
	}
}

export default Analytic;
