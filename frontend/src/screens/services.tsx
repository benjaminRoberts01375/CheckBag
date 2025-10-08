import "../styles.css";
import ServicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import { useState, useEffect, useRef } from "react";
import Service from "../types/service.tsx";
import { CgMoreVerticalAlt } from "react-icons/cg";
import { CommunicationProtocols, CommunicationProtocol } from "../types/strings";

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
					{service.external_address.map(externalAddress => (
						<ServiceStatus address={externalAddress} key={service.clientID} />
					))}
				</div>
				{service.internal_address ? <ServiceStatus address={service.internal_address} /> : null}
				<button onClick={() => dialogRef.current?.showModal()}>
					<CgMoreVerticalAlt className="icon" id={ServicesStyles["menu-icon"]} />
				</button>
			</div>
		</div>
	);
};

interface ServiceURLProps {
	address: string;
}

const ServiceStatus = ({ address }: ServiceURLProps) => {
	return <p className={ServicesStyles["service-status"]}>{address}</p>;
};

interface EditServiceProps {
	service: Service | undefined;
	finish: () => void;
}

const EditService = ({ service, finish }: EditServiceProps) => {
	const { serviceAdd, serviceDelete, serviceUpdate } = useList();
	const [title, setTitle] = useState(service?.title ?? "");
	const [incomingAddresses, setIncomingAddress] = useState(service?.external_address ?? []);

	const [outgoingProtocol, setOutgoingProtocol] = useState<CommunicationProtocol>(
		service?.internal_address?.startsWith("http") ? "http" : "https",
	);
	const [outgoingDomain, setOutgoingDomain] = useState(
		service?.internal_address.split(":")[1]?.substring(2) ?? "", // ex. https://www.google.com - splits on ":" and removes `//`
	);
	const [outgoingPort, setOutgoingPort] = useState(service?.internal_address.split(":")[2] ?? "80");

	function submit(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		service == undefined ? createService() : updateService(service);
		finish();

		// Reset the form
		setOutgoingProtocol("http");
		setOutgoingDomain("");
		setOutgoingPort("80");
		setIncomingAddress([]);
		setTitle("");
	}

	function createService() {
		console.log("Adding service:", outgoingProtocol + "://" + outgoingDomain + ":" + outgoingPort);
		serviceAdd(
			new Service(
				outgoingProtocol + "://" + outgoingDomain + ":" + outgoingPort,
				incomingAddresses,
				title,
			),
		);
	}

	function updateService(service: Service) {
		service.title = title;
		service.external_address = incomingAddresses;
		service.internal_address = outgoingProtocol + "://" + outgoingDomain + ":" + outgoingPort;
		serviceUpdate(service);
	}

	function cancel(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		console.log("Cancelling");
		finish();
		// Reset the form
		setOutgoingProtocol("http");
		setOutgoingDomain("");
		setOutgoingPort("80");
		setIncomingAddress([]);
		setTitle("");
	}

	function deleteService(e: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
		e.preventDefault();
		console.log("Deleting");
		if (service) {
			serviceDelete(service.clientID);
		}
		finish();
		// Reset the form
		setOutgoingProtocol("http");
		setOutgoingDomain("");
		setOutgoingPort("80");
		setIncomingAddress([]);
		setTitle("");
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
				value={incomingAddresses}
				onChange={e => setIncomingAddress([e.target.value])}
				className={ServicesStyles["input"]}
			/>
			<h3>To:</h3>
			<select
				value={outgoingProtocol}
				onChange={e => setOutgoingProtocol(e.target.value as CommunicationProtocol)}
			>
				{CommunicationProtocols.map(protocol => (
					<option value={protocol} key={protocol}>
						{protocol}
					</option>
				))}
			</select>
			<input
				type="url"
				autoComplete="off"
				placeholder="Forward Domain"
				value={outgoingDomain}
				onChange={e => setOutgoingDomain(e.target.value)}
				className={ServicesStyles["input"]}
			/>
			<input
				type="number"
				autoComplete="off"
				placeholder="Forward Port"
				value={outgoingPort}
				onChange={e => {
					var newPort = Number(e.target.value);
					if (newPort > 0 && newPort <= 65535) {
						setOutgoingPort(e.target.value);
					}
				}}
				min={1}
				max={65535}
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
