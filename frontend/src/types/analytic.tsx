class Analytic {
	quantity: number;
	country: Map<string, number>;
	ip: Map<string, number>;
	resource: Map<string, number>;
	responseCode: Map<number, number>;

	constructor(
		quantity: number,
		country: Map<string, number>,
		ip: Map<string, number>,
		resource: Map<string, number>,
		responseCode: Map<number, number>,
	) {
		this.quantity = quantity;
		this.country = country;
		this.ip = ip;
		this.resource = resource;
		this.responseCode = responseCode;
	}

	static fromJSON(data: any): Analytic {
		return new Analytic(
			data.quantity,
			new Map<string, number>(Object.entries(data.country)),
			new Map<string, number>(Object.entries(data.ip)),
			new Map<string, number>(Object.entries(data.resource)),
			new Map<number, number>(
				Object.entries(data.response_code).map(([key, value]) => [
					parseInt(key, 10),
					value as number,
				]),
			),
		);
	}
}

export default Analytic;
