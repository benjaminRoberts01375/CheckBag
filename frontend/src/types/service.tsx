import Analytic from "./analytic";
import { v4 as uuidv4 } from "uuid";
import GraphData from "./graph-data";
import { Timescale } from "./strings";
import ServiceURL from "./service-url";

// Individual service's processed data for a specific timescale
export interface ServiceProcessedData {
	quantityData: GraphData;
	responseCodeData: Map<number, number>;
	countryCodeData: Map<string, number>;
	ipAddressData: Map<string, number>;
	resourceUsage: Map<string, number>;
	totalRequests: number;
}

class Service {
	outgoing_address: ServiceURL;
	incoming_addresses: string[];
	title: string;
	id: string;
	clientID: string;
	enabled: boolean;

	// Pre-processed chart data for each timescale - this is all we need!
	hourProcessed: ServiceProcessedData;
	dayProcessed: ServiceProcessedData;
	monthProcessed: ServiceProcessedData;
	yearProcessed: ServiceProcessedData;

	constructor(
		outgoing_address: ServiceURL,
		incoming_addresses: string[],
		title: string,
		id: string = "",
		enabled: boolean = true,
		// Raw analytics - only used temporarily to create processed data
		hour: Map<string, Analytic> = new Map<string, Analytic>(),
		day: Map<string, Analytic> = new Map<string, Analytic>(),
		month: Map<string, Analytic> = new Map<string, Analytic>(),
		year: Map<string, Analytic> = new Map<string, Analytic>(),
	) {
		// Handle incoming addresses similarly
		incoming_addresses.forEach((address, index) => {
			if (address.substring(0, 4) !== "http") {
				address = "http://" + address;
			}
			const incomingURL = new URL(address);
			incoming_addresses[index] = incomingURL.port
				? `${incomingURL.hostname}:${incomingURL.port}`
				: incomingURL.hostname;
		});

		this.outgoing_address = outgoing_address;
		this.incoming_addresses = incoming_addresses;
		this.title = title;
		this.enabled = enabled;
		this.id = id;
		this.clientID = uuidv4();

		// Pre-process all timescales and discard raw data
		this.hourProcessed = this.processTimescaleData("hour", hour); // Has a short circuit to avoid parsing if empty
		this.dayProcessed = this.processTimescaleData("day", day);
		this.monthProcessed = this.processTimescaleData("month", month);
		this.yearProcessed = this.processTimescaleData("year", year);
	}

	// Process incoming analytics data to create graph-ready data
	private processTimescaleData(
		timescale: Timescale,
		analyticsMap: Map<string, Analytic>,
	): ServiceProcessedData {
		const responseCodesCounter = new Map<number, number>();
		const countryCounter = new Map<string, number>();
		const ipCounter = new Map<string, number>();
		const resourceCounter = new Map<string, number>();

		// Configure time span-specific settings
		let timeStepQuantity = 0;
		let rollback: (step: number) => Date;

		switch (timescale) {
			case "hour":
				rollback = (step: number) => {
					const now = new Date();
					const currentMinute = now.getMinutes();
					now.setUTCMinutes(currentMinute + step, 0, 0);
					return now;
				};
				timeStepQuantity = 60;
				break;
			case "day":
				rollback = (step: number) => {
					const now = new Date();
					now.setUTCHours(now.getUTCHours() + step, 0, 0, 0);
					return now;
				};
				timeStepQuantity = 24;
				break;
			case "month":
				rollback = (step: number): Date => {
					const now = new Date();
					now.setUTCHours(0, 0, 0, 0);
					const currentUTCDay = now.getUTCDate();
					now.setUTCDate(currentUTCDay + step);
					return now;
				};
				timeStepQuantity = 30;
				break;
			case "year":
				rollback = (step: number) => {
					const now = new Date();
					now.setUTCDate(1);
					now.setUTCHours(0, 0, 0, 0);
					const currentMonth = now.getUTCMonth();
					now.setUTCMonth(currentMonth + step);
					return now;
				};
				timeStepQuantity = 12;
				break;
		}

		const quantityData = new GraphData(this.title, timescale);
		var totalRequests = 0;
		// Generate data points for the entire time range
		for (let i = 0; i < timeStepQuantity; i++) {
			const date = rollback(-i);
			const dateString = () => {
				switch (timescale) {
					case "hour":
						return date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
					case "day":
						return date.toLocaleTimeString([], { hour: "numeric", hour12: true });
					case "month":
						return date.toLocaleString("default", {
							month: "short",
							day: "numeric",
						});
					case "year":
						return date.toLocaleString("default", {
							month: "short",
						});
					default:
						return date.toLocaleDateString();
				}
			};
			quantityData.x_values[timeStepQuantity - i - 1] = dateString();
			let requestQuantity = analyticsMap.get(date.toISOString())?.quantity ?? 0;
			quantityData.data[timeStepQuantity - i - 1] = requestQuantity;
			totalRequests += requestQuantity;

			const analytic = analyticsMap.get(date.toISOString());
			if (analytic !== undefined) {
				// Aggregate response codes
				analytic.responseCode.forEach((value, key) => {
					responseCodesCounter.set(key, (responseCodesCounter.get(key) ?? 0) + value);
				});
				// Aggregate countries
				analytic.country.forEach((value, key) => {
					countryCounter.set(key, (countryCounter.get(key) ?? 0) + value);
				});
				// Aggregate IP addresses
				analytic.ip.forEach((value, key) => {
					ipCounter.set(key, (ipCounter.get(key) ?? 0) + value);
				});
				// Aggregate resources
				analytic.resource.forEach((value, key) => {
					resourceCounter.set(key, (resourceCounter.get(key) ?? 0) + value);
				});
			}
		}

		return {
			quantityData,
			responseCodeData: responseCodesCounter,
			countryCodeData: countryCounter,
			ipAddressData: ipCounter,
			resourceUsage: resourceCounter,
			totalRequests: totalRequests,
		};
	}

	// Method to get processed data for a specific timescale
	getProcessedData(timescale: Timescale): ServiceProcessedData {
		switch (timescale) {
			case "hour":
				return this.hourProcessed;
			case "day":
				return this.dayProcessed;
			case "month":
				return this.monthProcessed;
			case "year":
				return this.yearProcessed;
			default:
				return this.hourProcessed;
		}
	}

	// Convert to JSON for server communication (only service config, no analytics)
	toJSON() {
		return {
			outgoing_address: this.outgoing_address,
			incoming_addresses: this.incoming_addresses,
			title: this.title,
			id: this.id,
		};
	}
	static fromJSON(data: any): Service {
		return new Service(
			new ServiceURL(
				data.outgoing_address.protocol,
				data.outgoing_address.domain,
				data.outgoing_address.port,
			),
			data.incoming_addresses,
			data.title,
			data.id || "",
			true,
			parseAnalyticMap(data.hour), // Has a short circuit to avoid parsing if empty
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
