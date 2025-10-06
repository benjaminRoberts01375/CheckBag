import { CommunicationProtocols } from "./strings";

class ServiceURL {
	protocol: CommunicationProtocols;
	hostname: string;
	port: number;

	constructor(protocol: CommunicationProtocols, hostname: string, port: number) {
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
