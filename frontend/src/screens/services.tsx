import "../styles.css";
import ServicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import { useState, useEffect, useRef } from "react";
import Service from "../types/service.tsx";
import { CgMoreVerticalAlt } from "react-icons/cg";
import ServiceURL from "../types/service-url.tsx";
import { CommunicationProtocols } from "../types/strings";

const ServicesScreen = () => {
	const { services, requestServiceData } = useList();
	useEffect(() => {
		requestServiceData();
	}, []);
	const dialogRef = useRef<HTMLDialogElement | null>(null);

	return (
		<div id={DashboardStyles["container"]}>
			<title>CheckBag - Services</title>
			<dialog ref={dialogRef}>
				<EditService service={undefined} finish={() => dialogRef.current?.close()} />
			</dialog>
			<div className={DashboardStyles["graph-group"]}>
				<h2 className="header">Services</h2>
				{services.map(service => (
					<ServiceEntry servicePass={service} key={service.clientID} />
				))}
				<button
					className="submit"
					id={ServicesStyles["add-service-button"]}
					onClick={() => dialogRef.current?.showModal()}
				>
					Add Service
				</button>
			</div>
		</div>
	);
};

export default ServicesScreen;

interface ServiceListEntryProps {
	servicePass: Service;
}

const ServiceEntry = ({ servicePass }: ServiceListEntryProps) => {
	const dialogRef = useRef<HTMLDialogElement | null>(null);
	const [service, _] = useState<Service>(servicePass);

	return (
		<div id={ServicesStyles["service-container"]}>
			<dialog ref={dialogRef}>
				<EditService service={service} finish={() => dialogRef.current?.close()} />
			</dialog>
			<h2>{service.title}</h2>
			<div className={ServicesStyles["connection-info"]}>
				<div id={ServicesStyles["service-endpoints"]}>
					{service.incoming_addresses.map(externalAddress => (
						<ServiceStatus address={externalAddress} key={service.clientID} />
					))}
				</div>
				{service.outgoing_address ? <ServiceStatus address={service.outgoing_address} /> : null}
				<button onClick={() => dialogRef.current?.showModal()}>
					<CgMoreVerticalAlt className="icon" id={ServicesStyles["menu-icon"]} />
				</button>
			</div>
		</div>
	);
};

interface ServiceURLProps {
	address: ServiceURL;
}

const ServiceStatus = ({ address }: ServiceURLProps) => {
	return <p className={ServicesStyles["service-status"]}>{address.toString()}</p>;
};

interface EditServiceProps {
	service: Service | undefined;
	finish: () => void;
}

const EditService = ({ service, finish }: EditServiceProps) => {
	const { serviceAdd, serviceDelete, serviceUpdate } = useList();
	const [title, setTitle] = useState(service?.title ?? "");

	const [outgoingAddressProtocol, setOutgoingAddressProtocol] = useState(
		service?.outgoing_address?.protocol ?? "http",
	);
	const [outgoingAddressDomain, setOutgoingAddressDomain] = useState(
		service?.outgoing_address?.hostname ?? "",
	);
	const [outgoingAddressPort, setOutgoingAddressPort] = useState(
		service?.outgoing_address?.port ?? 80,
	);

	const [incomingAddress, setIncomingAddress] = useState<string>("");

	function submit(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		service == undefined ? createService() : updateService(service);
		finish();
	}

	function createService() {
		serviceAdd(new Service(outgoingAddress, incomingAddresses, title));
	}

	function updateService(service: Service) {
		service.title = title;
		service.incoming_addresses = incomingAddresses;
		service.outgoing_address = outgoingAddress;
		serviceUpdate(service);
	}

	function cancel(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		console.log("Cancelling");
		finish();
		setTitle(service?.title ?? "");
		setIncomingAddress("");
		setOutgoingAddressProtocol(service?.outgoing_address?.protocol ?? "http");
		setOutgoingAddressDomain(service?.outgoing_address?.hostname ?? "");
		setOutgoingAddressPort(service?.outgoing_address?.port ?? 80);
	}

	function deleteService(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		console.log("Deleting");
		if (service) {
			serviceDelete(service.clientID);
		}
		finish();
	}

	return (
		<form id={ServicesStyles["edit-service-container"]}>
			<div id={ServicesStyles["edit-service-header"]}>
				<h1>Editing</h1>
				<input
					type="text"
					placeholder="Service Name"
					value={title}
					onChange={e => setTitle(e.target.value)}
					className={ServicesStyles["input"]}
				/>
			</div>
			<h3>From:</h3>
			<input
				type="url"
				autoComplete="off"
				placeholder="Forward Address"
				value={incomingAddress}
				onChange={e => setIncomingAddress(e.target.value)}
				className={ServicesStyles["input"]}
			/>
			<h3>To:</h3>
			<select value={outgoingAddressProtocol}>
				{CommunicationProtocols.map(protocol => (
					<option value={protocol}>{protocol}</option>
				))}
			</select>
			<input
				type="url"
				autoComplete="off"
				placeholder="Domain"
				value={outgoingAddressDomain}
				onChange={e => setOutgoingAddressDomain(e.target.value)}
				className={ServicesStyles["input"]}
			/>
			<input
				type="number"
				autoComplete="off"
				placeholder="Port"
				value={outgoingAddressPort}
				onChange={e => setOutgoingAddressPort(parseInt(e.target.value))}
				className={ServicesStyles["input"]}
			/>
			<div id={ServicesStyles["submission-buttons"]}>
				{service ? (
					<button
						className="delete"
						role="delete"
						onClick={e => {
							deleteService(e);
						}}
					>
						Delete
					</button>
				) : null}

				<button
					className="cancel"
					role="cancel"
					onClick={e => {
						cancel(e);
					}}
				>
					Cancel
				</button>
				<button
					className="submit"
					role="submit"
					onClick={e => {
						submit(e);
					}}
				>
					{service == undefined ? "Create" : "Save"}
				</button>
			</div>
		</form>
	);
};
