import Analytic from "./analytic";

class Service {
	internal_address: string;
	external_address: string[];
	title: string;
	id: string;
	clientID: string;
	hour: Map<Date, Analytic>;
	day: Map<Date, Analytic>;
	month: Map<Date, Analytic>;
	year: Map<Date, Analytic>;
	enabled: boolean;

	constructor(
		internal_address: string,
		external_address: string[],
		title: string,
		enabled: boolean = false,
		id: string = "",
		hour: Map<Date, Analytic> = new Map<Date, Analytic>(),
		day: Map<Date, Analytic> = new Map<Date, Analytic>(),
		month: Map<Date, Analytic> = new Map<Date, Analytic>(),
		year: Map<Date, Analytic> = new Map<Date, Analytic>(),
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
}

export default Service;
