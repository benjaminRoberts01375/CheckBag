class ResourceUsageData {
	service: string;
	resource: string;
	quantity: number;

	constructor(service: string, resource: string, quantity: number) {
		this.service = service;
		this.resource = resource;
		this.quantity = quantity;
	}
}

export default ResourceUsageData;
