class ChartData {
	value: number;
	label: string;
	id: string | number;

	constructor(value: number, label: string) {
		this.value = value;
		this.label = label;
		this.id = label;
	}
}

export default ChartData;
