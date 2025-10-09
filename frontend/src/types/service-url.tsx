import { CommunicationProtocol } from "./strings";

class ServiceURL {
	protocol: CommunicationProtocol;
	hostname: string;
	port: number;

	constructor(protocol: CommunicationProtocol = "http", hostname: string = "", port: number = 80) {
		this.protocol = protocol;
		this.hostname = hostname;
		this.port = port;
		if (this.port === 0) {
			this.port = this.protocol === "http" ? 80 : 443;
		}
	}
	toString(): string {
		return `${this.protocol}://${this.hostname}:${this.port}`;
	}
}

export default ServiceURL;
