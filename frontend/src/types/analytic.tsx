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

	static fromJSON(data: any): Analytic {
		return new Analytic(
			data.quantity,
			new Map<string, number>(Object.entries(data.country || {})),
			new Map<string, number>(Object.entries(data.ip || {})),
			new Map<string, number>(Object.entries(data.resource || {})),
		);
	}
}

export default Analytic;
