class ResourceUsageData {
	service: string;
	resource: string;
	quantity: number;
	service_url: string;

	constructor(service_name: string, service_url: string, resource: string, quantity: number) {
		this.service = service_name;
		this.resource = resource;
		this.quantity = quantity;
		this.service_url = service_url;
	}
}

export default ResourceUsageData;
