import { Timescale } from "./strings";

class GraphData {
	data: number[]; // Quantity per hour per service. Index 0 = left, index 60/24/30/12 = right
	label: string; // Service title
	stack = "total" as const;
	type = "bar" as const;
	x_values: string[];

	constructor(label: string, timescale: Timescale) {
		this.label = label;
		var indices = 60;
		switch (timescale) {
			case "day":
				indices = 24;
				break;
			case "month":
				indices = 30;
				break;
			case "year":
				indices = 12;
				break;
		}
		this.data = new Array(indices).fill(0);
		this.x_values = new Array(indices).fill("");
	}
}

export default GraphData;
