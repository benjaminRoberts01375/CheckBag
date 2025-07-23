import Analytic from "./analytic";

class Service {
	internal_address: string;
	external_address: string[];
	title: string;
	id: string;
	clientID: string;
	hour: Map<string, Analytic>;
	day: Map<string, Analytic>;
	month: Map<string, Analytic>;
	year: Map<string, Analytic>;
	enabled: boolean;

	constructor(
		internal_address: string,
		external_address: string[],
		title: string,
		id: string = "",
		hour: Map<string, Analytic> = new Map<string, Analytic>(),
		day: Map<string, Analytic> = new Map<string, Analytic>(),
		month: Map<string, Analytic> = new Map<string, Analytic>(),
		year: Map<string, Analytic> = new Map<string, Analytic>(),
		enabled: boolean = true,
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
		this.enabled = enabled;
		this.id = id;
		this.clientID = crypto.randomUUID();
		this.hour = hour;
		this.day = day;
		this.month = month;
		this.year = year;
	}

	static fromJSON(data: any): Service {
		return new Service(
			data.internal_address,
			data.external_address,
			data.title,
			data.id || "",
			parseAnalyticMap(data.hour),
			parseAnalyticMap(data.day),
			parseAnalyticMap(data.month),
			parseAnalyticMap(data.year),
		);
	}
}

function parseAnalyticMap(data: any): Map<string, Analytic> {
	const analyticsMap = new Map<string, Analytic>();
	if (data) {
		for (const key in data) {
			if (data.hasOwnProperty(key)) {
				analyticsMap.set(new Date(key).toISOString(), Analytic.fromJSON(data[key]));
			}
		}
	}
	return analyticsMap;
}

export default Service;
