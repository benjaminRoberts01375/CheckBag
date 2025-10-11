import { CommunicationProtocol } from "./strings";

class ServiceURL {
	protocol: CommunicationProtocol;
	domain: string;
	port: number;

	constructor(protocol: CommunicationProtocol = "http", domain: string = "", port: number = 80) {
		this.protocol = protocol;
		this.domain = domain;
		this.port = port;
		if (this.port === 0) {
			this.port = this.protocol === "http" ? 80 : 443;
		}
	}
	public toString(): string {
		return `${this.protocol}://${this.domain}:${this.port}`;
	}
}

export default ServiceURL;
