import { Timescale, Timescales } from "./types/strings";

self.onmessage = async (e: MessageEvent<{ url: string }>) => {
	// Process the data
	const processedData = await requestData(e.data.url, [Timescales[0]]);

	// Send back with timeframe identifier
	self.postMessage({ timeframe: Timescales[0], data: processedData });
};

async function requestData(url: string, timeSteps: Timescale[]) {
	const requestURL = new URL(url);
	var rawJSON: any;

	try {
		timeSteps.forEach(timeStep => {
			requestURL.searchParams.append("time-step", timeStep);
		});
		const response = await fetch(requestURL.toString(), {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
			},
			credentials: "include",
		});

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}

		rawJSON = await response.json();
	} catch (error) {
		console.error("Error requesting data:", error);
		throw error;
	}

	console.log(rawJSON);
}
